package useraction

import (
	"github.com/techfront/core/src/kernel/auth"
	"github.com/techfront/core/src/kernel/router"
)

/**
* HandleLogout очищает сессию для авторизованного пользователя
 */
func HandleLogout(context router.Context) error {

	auth.ClearSession(context.Writer())

	// Редирект на главную
	return router.Redirect(context, "/")
}
