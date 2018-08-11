package commentaction

import (
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/component/comment"

	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleUpdateShow responds to GET /comments/update with the form to update a comment
 */
func HandleUpdateShow(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение комментария
	commentEntity, err := comment.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Проверка прав на действие
	if err = authorise.Resource(context, commentEntity); err != nil {
		return router.NotAuthorizedError(err, displayerror.AccessDenied...)
	}

	// Отображение шаблона
	v := view.New(
		"component/comment/template/update",
		"component/comment/template/form",
		"component/comment/template/form_embed",
	)
	v.Vars["comment"] = commentEntity

	return v.Render(context)
}

/**
* HandleUpdate responds to POST /comments/update
 */
func HandleUpdate(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение комментария
	commentEntity, err := comment.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Проверка прав
	err = authorise.ResourceAndAuthenticity(context, commentEntity)
	if err != nil {
		return router.NotAuthorizedError(err, displayerror.AccessDenied...)
	}

	// Очистка параметров для текущей роли
	accepted := comment.AllowedParams()
	if authorise.CurrentUser(context).Admin() {
		accepted = comment.AllowedParamsAdmin()
	}
	cleanedParams := params.Clean(accepted)

	// Обновление комментария
	err = commentEntity.Update(cleanedParams)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Редирект к топику
	return router.Redirect(context, commentEntity.URLTopic())
}
