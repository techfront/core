package topicaction

import (
	"fmt"
	"strings"

	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/cache"
	"github.com/techfront/core/src/lib/displayerror"

	"github.com/techfront/core/src/component/comment"
	"github.com/techfront/core/src/component/topic"
)

/**
* HandleShow displays a single topic.
 */
func HandleShow(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение топика
	topicEntity, err := topic.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Скрываем топик если он отклонен и пользователь не админ.
	if topicEntity.IsRejected() && !authorise.CurrentUser(context).Admin() {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Редирект на канонический URL
	if context.Path() != topicEntity.URLShow() {
		return router.Redirect(context, topicEntity.URLShow())
	}

	// Получение интересных топиков для виджета
	// Оптимизация: кэшированние выборки
	relatedTopicKey := fmt.Sprintf("topic__id-%d--related", params.GetInt("id"))

	var relatedTopicList []*topic.Topic

	if err := cache.Get(relatedTopicKey, &relatedTopicList); err != nil {
		relatedTopicQuery := topicEntity.Related().Limit(8)

		results, err := topic.FindAll(relatedTopicQuery)
		if err != nil {
			return router.InternalError(err, displayerror.UnknownError...)
		}

		relatedTopicList = results

		if err := cache.Set(relatedTopicKey, relatedTopicList, 360); err != nil {
			return router.InternalError(err, displayerror.UnknownError...)
		}
	}

	// Получение комментариев
	commentsQuery := comment.Where("comment_id_topic=?", topicEntity.Id).Order(comment.RANK_ORDER)
	commentList, err := comment.FindAllWithChild(commentsQuery)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Отображаение шаблона
	v := view.New(
		"component/topic/template/show",
		"component/topic/template/row_mini",
		"component/topic/template/widget/widget_related",
		"component/topic/template/widget/widget_share",
		"component/comment/template/comment",
		"component/comment/template/form_embed",
	)

	switch params.Get("message") {
	case "info--topic_in_the_queue":
		v.Vars["message_type"] = "default"
		v.Vars["message"] = "Спасибо, топик успешно создан и ожидает свою очередь."
	case "error--comment_invalid_create_timeout":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Добавить новый комментарий можно через 1 минуту. Подождите и попробуйте снова."
	case "error--comment_invalid_length":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Комментарий должен быть длинее 10 символов."
	}

	var desc string
	if len(topicEntity.Text) > 155 {
		desc = topicEntity.Text[:155]
		lastWordLenght := len(desc) - strings.LastIndex(desc, " ")
		desc = fmt.Sprintf("%s...", desc[:155 - lastWordLenght])
	}

	v.Vars["topic"] = topicEntity
	v.Vars["meta_title"] = topicEntity.Name + " / Техфронт"
	v.Vars["meta_desc"] = desc
	v.Vars["meta_keywords"] = ""
	v.Vars["comments"] = commentList
	v.Vars["related_topics"] = relatedTopicList

	return v.Render(context)
}
