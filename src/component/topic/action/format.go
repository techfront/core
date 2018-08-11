package topicaction

import (
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/component/topic"
	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleFormat отображает список топиков на странице формата
 */
func HandleFormat(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Построение запроса
	query := topic.Published().Where("score(topic_count_upvote, topic_count_downvote, topic_count_flag) > -6")

	query.Order("rank(score(topic_count_upvote, topic_count_downvote, topic_count_flag), topic_created_at) desc, score(topic_count_upvote, topic_count_downvote, topic_count_flag) desc, topic_id desc").Limit(DEFAULT_LIMIT)

	var formatTitle string
	var formatID int64

	// Все форматы
	formatID = -1

	// Определение конкретного формата с помощью параметра в URL
	switch params.Get("format") {
	case "topic":
		formatTitle = "топики"
		formatID = topic.FormatTopic
		query.Where("topic_id_format = ? OR topic_id_format IS NULL", topic.FormatTopic)
	case "news":
		formatTitle = "новости"
		formatID = topic.FormatNews
		query.Where("topic_id_format = ?", topic.FormatNews)
	case "project":
		formatTitle = "проекты"
		formatID = topic.FormatProject
		query.Where("topic_id_format = ?", topic.FormatProject)
	case "podcast":
		formatTitle = "подкасты"
		formatID = topic.FormatPodcast
		query.Where("topic_id_format = ?", topic.FormatPodcast)
	case "video":
		formatTitle = "видео"
		formatID = topic.FormatVideo
		query.Where("topic_id_format = ?", topic.FormatVideo)
	case "question":
		formatTitle = "вопросы"
		formatID = topic.FormatQuestion
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
	v := view.New(
		"component/topic/template/index",
		"component/topic/template/row_mini",
	)

	v.Vars["page"] = currentPage
	v.Vars["topics"] = results
	v.Vars["topics_count"] = len(results)
	v.Vars["topics_format"] = formatID

	if err := setTopicMetadata(v, context); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	v.Vars["page_title"] = "Все " + formatTitle
	v.Vars["meta_title"] = "Все " + formatTitle + " / Техфронт"

	return v.Render(context)
}
