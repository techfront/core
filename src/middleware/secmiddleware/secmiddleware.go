package secmiddleware

import (
	"github.com/techfront/core/src/kernel/router"
)

/**
*  Secure Middleware. Посредник добавляет заголовки безопасности.
*
* @param next router.Handler следующий хендлер.
 */
func Middleware(next router.Handler) router.Handler {
	return router.HandlerFunc(func(context router.Context) error {
		if should(context) {
			w := context.Writer()
			w.Header().Set("X-Xss-Protection", "1; mode=block")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			context.SetWriter(w)
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
	r := context.Request()
	if r.URL.Path == "/away" {
		return false
	}

	ext := context.Ext()
	if len(ext) > 0 {
		return false
	}

	return true
}
