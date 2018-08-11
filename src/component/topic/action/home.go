package topicaction

import (
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/component/topic"
	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleHome displays a list of topics using gravity to order them
* used for the home page for gravity rank see votes.go
* responds to GET /
 */
func HandleHome(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Построение запроса
	query := topic.Published().Limit(DEFAULT_LIMIT)

	query.Where("topic_count_upvote > 0")

	query.Order("topic_status desc, rank(score(topic_count_upvote, topic_count_downvote, topic_count_flag), topic_created_at) desc, topic_id desc")

	currentPage := int(params.GetInt("page"))
	if currentPage > 0 {
		query.Offset(DEFAULT_LIMIT * currentPage)
	}

	// Получение запроса
	results, err := topic.FindAll(query)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Отображение шаблона
	v := view.New(
		"component/topic/template/home",
		"component/topic/template/row_mini",
	)

	v.Vars["topics"] = results

	if err := setTopicMetadata(v, context); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return v.Render(context)
}
