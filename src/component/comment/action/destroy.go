package commentaction

import (
	"github.com/techfront/core/src/kernel/router"

	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"

	"github.com/techfront/core/src/component/comment"
)

/**
* HandleDestroy handles a DESTROY request for comments
 */
func HandleDestroy(context router.Context) error {

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

	// Удаление комментария
	if err := commentEntity.Destroy(); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Редирект к топику
	return router.Redirect(context, commentEntity.URLTopic())
}
