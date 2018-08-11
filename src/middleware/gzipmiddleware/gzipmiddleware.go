package gzipmiddleware

import (
	"bufio"
	"github.com/klauspost/pgzip"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/gzip"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

type GzipWriter struct {
	http.ResponseWriter
	writer *pgzip.Writer
	size   int64
}

func (gw *GzipWriter) Write(data []byte) (int, error) {
	if gw.Header().Get("Content-Type") == "" {
		gw.Header().Set("Content-Type", http.DetectContentType(data))
	}

	n, err := gw.writer.Write(data)

	gw.size += int64(n)

	return n, err
}

func (gw *GzipWriter) WriteHeader(code int) {
	if code == http.StatusNoContent {
		gw.Header().Del("Content-Encoding")
	}

	gw.ResponseWriter.WriteHeader(code)
}

func (gw *GzipWriter) Flush() {
	gw.ResponseWriter.(http.Flusher).Flush()
}

func (gw *GzipWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return gw.ResponseWriter.(http.Hijacker).Hijack()
}

func (gw *GzipWriter) CloseNotify() <-chan bool {
	return gw.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

/**
* Посредник выполняет gzip-сжатие текущего ресурса
*
* @param next router.Handler следующий хендлер.
 */
func Middleware(next router.Handler) router.Handler {
	return router.HandlerFunc(func(context router.Context) error {

		if should(context) {
			w := context.Writer()

			gz := gzip.Pool.Get().(*pgzip.Writer)
			defer gzip.Pool.Put(gz)
			gz.Reset(w)

			// Добавление заголовков ответа
			w.Header().Add("Content-Encoding", "gzip")
			w.Header().Add("Vary", "Accept-Encoding")

			// Переопределение контекста
			gw := &GzipWriter{ResponseWriter: w, writer: gz, size: 0}
			context.SetWriter(gw)

			defer func() {
				if gw.size == 0 {
					w.Header().Del("Content-Encoding")
					context.SetWriter(w)
					gz.Reset(ioutil.Discard)
				}

				gz.Close()
			}()

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

	// Если браузер не поддерживает Gzip, то False
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		return false
	}

	// Если редирект, то False
	if r.URL.Path == "/away" {
		return false
	}

	// Если изображение, то False
	ext := context.Ext()
	switch ext {
	case ".png", ".gif", ".jpeg", ".jpg":
		return false
	}

	// Запросы к API - False
	if strings.HasPrefix(r.URL.Path, "/api") {
		return false
	}

	return true
}
