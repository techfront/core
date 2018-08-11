package ascemiddleware

import (
	"bytes"
	"fmt"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/cache"
	"net/http"
)

/**
* Тип CachedWriter это структура для буфферизации http.ResponseWriter
 */
type CachedWriter struct {
	http.ResponseWriter
	buffer bytes.Buffer
}

/**
* Тип Response это структура закэшированного ответа.
 */
type Response struct {
	Status int
	Header http.Header
	Body   []byte
}

func (cw *CachedWriter) Write(data []byte) (int, error) {
	n, err := cw.ResponseWriter.Write(data)
	if err == nil {
		cw.buffer.Write(data)
	}

	return n, err
}

/**
* Посредник кеширует Get-запрос для анонимуса.
*
* @param next router.Handler следующий хендлер.
 */
func Middleware(next router.Handler) router.Handler {
	return router.HandlerFunc(func(context router.Context) error {

		if should(context) {

			// Получение ключа
			hash := cache.ContextToHash(context)
			key := fmt.Sprintf("cache:asce:%d", hash)

			var responseCache Response

			if err := cache.Get(key, &responseCache); err == nil {

				// Добавляем заголовки
				for k, vals := range responseCache.Header {
					for _, v := range vals {
						context.Writer().Header().Set(k, v)
					}
				}

				// Добавляем параметр X-Cache
				context.Writer().Header().Set("X-Cache", "HIT")

				// И возвращаем ответ
				context.Writer().Write(responseCache.Body)

				return nil
			}

			// Получение Writter
			w := context.Writer()

			// Переопределение Writer на кэшируемый
			cw := &CachedWriter{ResponseWriter: w}
			newContext := context
			newContext.SetWriter(cw)

			newContext.Writer().Header().Set("X-Cache", "MISS")

			if err := next.ServeHTTP(newContext); err != nil {
				return err
			}

			defer func() {
				context.SetWriter(w)
			}()

			// Формируем ответ
			response := Response{
				Status: 200,
				Header: http.Header(cw.Header()),
				Body:   cw.buffer.Bytes(),
			}

			// Записываем ответ в кэш
			if err := cache.Set(key, response, 180); err != nil {
				return nil
			}

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
	r := context.Request()
	if r.Method != "GET" {
		return false
	}

	if r.URL.Path == "/away" {
		return false
	}

	// Проверка расширения ресурса
	ext := context.Ext()
	if len(ext) > 0 {
		return false
	}

	// Если пользователь анонимус
	if authorise.CurrentUser(context).Id != 0 {
		return false
	}

	return true
}
