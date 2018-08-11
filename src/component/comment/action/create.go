package commentaction

import (
	"fmt"

	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/component/comment"
	"github.com/techfront/core/src/component/topic"

	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleCreateShow serves the create form via GET for comments
 */
func HandleCreateShow(context router.Context) error {

	// Проверка прав
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err, displayerror.AccessDenied...)
	}

	// Отображение шаблона
	v := view.New()

	commentEntity := comment.New()

	v.Vars["comment"] = commentEntity

	return v.Render(context)
}

/**
* HandleCreate handles the POST of the create form for comments
 */
func HandleCreate(context router.Context) error {

	// Проверка прав
	err := authorise.AuthenticityToken(context)
	if err != nil {
		return router.NotAuthorizedError(err, displayerror.AccessDenied...)
	}

	if !authorise.CurrentUser(context).CanComment() {
		return router.NotAuthorizedError(nil, displayerror.AccessDeniedKarma...)
	}

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Поиск топика-родителя по отправленном ID
	topicEntity, err := topic.Find(params.GetInt("comment_id_topic"))
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	text := params.Get("comment_text")
	if len(text) < 10 {
		return router.Redirect(context, fmt.Sprintf("/topics/%d?message=error--comment_invalid_length", topicEntity.Id))
	}

	params.SetInt("comment_id_topic", topicEntity.Id)

	// Поиск пользователя
	userEntity := authorise.CurrentUser(context)
	ip := context.ClientIP()

	// Проверка лимита
	limitValid, err := userEntity.CheckCreateCommentLimit()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}
	if !limitValid {
		return router.Redirect(context, fmt.Sprintf("/topics/%d?message=error--comment_invalid_create_timeout", topicEntity.Id))
	}

	// Очистка параметров
	accepted := comment.AllowedParams()
	cleanedParams := params.Clean(accepted)

	// Определение базовых параметров
	cleanedParams["comment_status"] = "100"
	cleanedParams["comment_id_user"] = fmt.Sprintf("%d", userEntity.Id)
	cleanedParams["comment_count_upvote"] = "1"

	// Определение положения комментария
	parentID := params.GetInt("comment_id_parent")
	if parentID > 0 {
		parent, err := comment.Find(parentID)
		if err != nil {
			return router.InternalError(err, displayerror.UnknownError...)
		}
		cleanedParams["comment_dotted_ids"] = fmt.Sprintf(parent.DottedIds + ".")
	}

	// Создание комментариев
	id, err := comment.Create(cleanedParams)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Обновление колличества комментариев пользователя
	userParams := map[string]string{"user_count_comment": fmt.Sprintf("%d", userEntity.CommentCount+1)}
	err = userEntity.Update(userParams)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Обновление колличества комментариев топика
	topicParams := map[string]string{"topic_count_comment": fmt.Sprintf("%d", topicEntity.CommentCount+1)}
	err = topicEntity.Update(topicParams)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Поиск созданного комментария
	commentEntity, err := comment.Find(id)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Добавляем первый голос
	err = commentRecordUpvote(commentEntity, userEntity, ip)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return router.Redirect(context, commentEntity.URLTopic())
}
