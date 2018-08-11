package useraction

import (
	"github.com/techfront/core/src/component/user"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleDestroy responds to POST /users/1/destroy
 */
func HandleDestroy(context router.Context) error {

	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение пользователя
	userEntity, err := user.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Проверка прав
	err = authorise.ResourceAndAuthenticity(context, userEntity)
	if err != nil {
		return router.NotAuthorizedError(err, displayerror.AccessDenied...)
	}

	// Удаление пользователя
	err = userEntity.Destroy()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Редирект на главную
	return router.Redirect(context, "/")
}
