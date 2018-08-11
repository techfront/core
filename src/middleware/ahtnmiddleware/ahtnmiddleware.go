package ahtnmiddleware

import (
	"github.com/techfront/core/src/kernel/auth"
	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/kernel/router"
)

/**
 * Посредник добавляет в контекст токен для csfr.
 *
 * @param next router.Handler следующий хендлер.
 */
func Middleware(next router.Handler) router.Handler {
	return router.HandlerFunc(func(context router.Context) error {
		if should(context) {	
			token, err := auth.AuthenticityToken(context.Writer(), context.Request())
			if err != nil {
				return err
			}
			context.Set("authenticity_token", token)
		}

		return next.ServeHTTP(context)
	})
}

/**
 * Функция should проверяет должна ли работать Middleware для данного запроса.
 *
 * @param context Контекст текущего ресурса
 * @return bool
 */
func should(context router.Context) bool {
	ext := context.Ext()
	if len(ext) > 0 {
		return false
	}

	if authorise.CurrentUser(context).Id == 0 {
		return false
	}

	return true
}