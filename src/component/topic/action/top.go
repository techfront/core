package topicaction

import (
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	// "github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"

	"github.com/techfront/core/src/component/topic"
)

/**
* HandleTop displays a list of topics the user has top in the past
 */
func HandleTop(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Построение запроса
	query := topic.Published().Limit(DEFAULT_LIMIT)

	query.Where("score(topic_count_upvote, topic_count_downvote, topic_count_flag) > 0")
	query.Order("rank(score(topic_count_upvote, topic_count_downvote, topic_count_flag), topic_created_at) desc, score(topic_count_upvote, topic_count_downvote, topic_count_flag) desc, topic_id desc")

	/*
		// Выбрать для авторизованного пользователя
		userEntity := authorise.CurrentUser(context)
		if !user.Anon() {
			// Can we use a join instead?
			v := query.New("votes", "topic_id").Select("select topic_id as id from votes").Where("user_id=? AND topic_id IS NOT NULL AND points > 0", u.Id)

			topicIDs := v.ResultIDs()
			if len(topicIDs) > 0 {
				query.WhereIn("id", topicIDs)
			}
		}
	*/

	// Определение текущей страницы и добавление отступа
	currentPage := int(params.GetInt("page"))
	if currentPage > 0 {
		query.Offset(DEFAULT_LIMIT * currentPage)
	}

	// Получение топиков
	results, err := topic.FindAll(query)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Отображение шаблона
	v := view.New(
		"component/topic/template/index",
		"component/topic/template/row_mini",
	)

	v.Vars["page"] = currentPage
	v.Vars["topics"] = results
	v.Vars["topics_format"] = -1
	v.Vars["topics_count"] = len(results)

	if params.Get("format") == ".xml" {
		v.Extension = "xml.got"
		v.Templates = []string{"component/topic/template/index"}
	}

	if err := setTopicMetadata(v, context); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	v.Vars["meta_title"] =  "Топ топиков / Техфронт"

	return v.Render(context)
}
