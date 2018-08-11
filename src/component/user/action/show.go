package useraction

import (
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/component/comment"
	"github.com/techfront/core/src/component/topic"
	"github.com/techfront/core/src/component/user"

	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleShow serve a get request at /users/1
 */
func HandleShow(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Поиск пользователя
	userEntity, err := user.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Получение комментариев пользователя
	userCommentQuery := comment.Where("comment_id_user=?", userEntity.Id).Limit(24).Order("comment_created_at desc")
	userCommentList, err := comment.FindAllWithChild(userCommentQuery)
	if err != nil {
		return router.InternalError(nil, displayerror.UnknownError...)
	}

	// Получение топиков пользователя
	userTopicQuery := topic.Published().Where("topic_id_user=?", userEntity.Id).Limit(12).Order("topic_created_at desc")
	userTopicList, err := topic.FindAll(userTopicQuery)
	if err != nil {
		return router.InternalError(nil, displayerror.UnknownError...)
	}

	// Отображение шаблона
	v := view.New(
		"component/user/template/show",
		"component/comment/template/comment",
		"component/comment/template/form_embed",
		"component/topic/template/row_mini",
	)

	v.Vars["user"] = userEntity
	v.Vars["comments"] = userCommentList
	v.Vars["topics"] = userTopicList
	v.Vars["meta_title"] = userEntity.Name + " / Техфронт"

	return v.Render(context)
}
