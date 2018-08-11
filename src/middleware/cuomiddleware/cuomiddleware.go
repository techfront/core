package cuomiddleware

import (
	"fmt"
	"github.com/fragmenta/query"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/authorise"
	"time"

	"github.com/techfront/core/src/lib/cache"
	"github.com/techfront/core/src/lib/displayerror"
)

type Response struct {
	Id        int64
	VisitedAt time.Time
}

/**
* Посредник фиксирует статус текущего пользователя.
*
* @param next router.Handler следующий обработчик.
 */
func Middleware(next router.Handler) router.Handler {
	return router.HandlerFunc(func(context router.Context) error {
		if should(context) {
			currentUser := authorise.CurrentUser(context)
			currentTime := time.Now().UTC()
			key := fmt.Sprintf("cache:cuo:%d", currentUser.Id)
			var respCache Response

			if err := cache.Get(key, &respCache); err != nil {
				if time.Since(respCache.VisitedAt).Seconds() > 360 {
					currentUserParams := map[string]string{"user_visited_at": query.TimeString(currentTime)}

					err = currentUser.Update(currentUserParams)
					if err != nil {
						return router.InternalError(err, displayerror.UnknownError...)
					}

					resp := Response{
						Id:        currentUser.Id,
						VisitedAt: currentTime,
					}

					cache.Set(key, resp, 360)
				}
			}
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

	// Если пользователь анонимус
	if authorise.CurrentUser(context).Id == 0 {
		return false
	}

	return true
}
