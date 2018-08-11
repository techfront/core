package logmiddleware

import (
	"github.com/techfront/core/src/kernel/router"
	"time"
)

/**
* Посредник логирует текущий запрос.
*
* @param next router.Handler следующий хендлер.
 */
func Middleware(next router.Handler) router.Handler {
	return router.HandlerFunc(func(context router.Context) error {
		if should(context) {
			started := time.Now()
			path := context.Path()
			request := context.Request()
			ip := context.ClientIP()

			context.Logf("#info Started %s %s for %s", request.Method, path, ip)

			err := next.ServeHTTP(context)
			if err != nil {
				return err
			}

			defer func() {
				end := time.Since(started).String()
				status := 200
				context.Logf("#info Finished %s %s for %s status %d in %s", request.Method, path, ip, status, end)
			}()

			return nil
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

	return true
}
