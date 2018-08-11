package topicaction

import (
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"
	"time"

	"github.com/techfront/core/src/component/topic"
	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleFormatXml отображает список топиков на странице формата в XML
 */
func HandleFormatXml(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Построение запроса
	query := topic.Published().Where("score(topic_count_upvote, topic_count_downvote, topic_count_flag) > -6")
	query.Order("rank(score(topic_count_upvote, topic_count_downvote, topic_count_flag), topic_created_at) desc, score(topic_count_upvote, topic_count_downvote, topic_count_flag) desc, topic_id desc").Limit(DEFAULT_LIMIT)

	var formatTitle string

	switch params.Get("format") {
	case "topics":
		formatTitle = "топики"
		query.Where("topic_id_format = ? OR topic_id_format IS NULL", topic.FormatTopic)
	case "news":
		formatTitle = "новости"
		query.Where("topic_id_format = ?", topic.FormatNews)
	case "projects":
		formatTitle = "проекты"
		query.Where("topic_id_format = ?", topic.FormatProject)
	case "podcasts":
		formatTitle = "подкасты"
		query.Where("topic_id_format = ?", topic.FormatPodcast)
	case "video":
		formatTitle = "видео"
		query.Where("topic_id_format = ?", topic.FormatVideo)
	case "questions":
		formatTitle = "вопросы"
		query.Where("topic_id_format = ?", topic.FormatQuestion)
	default:
		return router.NotFoundError(nil, displayerror.PageNotFound...)
	}

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
	v := view.New("component/topic/template/index")

	v.Extension = "xml.got"
	v.Layout = "index"
	v.Format = "application/rss+xml"
	v.Templates = []string{"component/topic/template/index"}

	v.Vars["page"] = currentPage
	v.Vars["topics"] = results
	v.Vars["topics_count"] = len(results)
	v.Vars["meta_title"] = "Все " + formatTitle + " / Техфронт"
	v.Vars["pub_date"] = time.Now()

	return v.Render(context)
}
