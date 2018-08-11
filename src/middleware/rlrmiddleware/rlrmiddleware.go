package rlrmiddleware

import (
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/limiter"
	"log"
	"net/http"
)

/**
* Посредник ограничивает колличество запросов к хендлеру.
*
* @param next router.Handler следующий хендлер.
 */
func Middleware(next router.Handler) router.Handler {
	return router.HandlerFunc(func(context router.Context) error {
		if should(context) {
			nh := wrap(next, context)

			l, err := limiter.New(100, "minute", 0)
			if err != nil {
				return err
			}

			lh := l.ByPath()
			lh.RateLimit(nh).ServeHTTP(context.Writer(), context.Request())

			return nil

		}

		return next.ServeHTTP(context)
	})
}

func wrap(handler router.Handler, context router.Context) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		err := handler.ServeHTTP(context)
		if err != nil {
			log.Print(err)
		}

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
