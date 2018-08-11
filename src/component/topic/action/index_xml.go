package topicaction

import (
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"
	"time"

	"github.com/techfront/core/src/component/topic"
	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleIndexXml displays a list of topics at /xml/topics
 */
func HandleIndexXml(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Построение запроса
	query := topic.Published().Limit(DEFAULT_LIMIT)

	query.Where("score(topic_count_upvote, topic_count_downvote, topic_count_flag) > -6").Order("topic_created_at desc")

	// Определяем текущую страницу и делаем отступ
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
	v := view.New("component/topic/template/index")

	v.Extension = "xml.got"
	v.Layout = "index"
	v.Format = "application/rss+xml"
	v.Templates = []string{"component/topic/template/index"}

	v.Vars["page"] = currentPage
	v.Vars["topics"] = results
	v.Vars["topics_count"] = len(results)
	v.Vars["pub_date"] = time.Now()

	return v.Render(context)
}
