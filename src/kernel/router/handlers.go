package router

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

func fileHandler(context Context) error {

	localPath := "./public" + path.Clean(context.Path())

	if _, err := os.Stat(localPath); err != nil {
		if os.IsNotExist(err) {
			return NotFoundError(err)
		}

		return NotAuthorizedError(err)
	}

	http.ServeFile(context.Writer(), context.Request(), localPath)

	return nil
}

func errorHandler(context Context, e error) {

	err := ToStatusError(e)

	writer := context.Writer()

	// Set the headers
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(err.Status)

	// Write a simple error message page
	html := fmt.Sprintf("<h1>%s</h1><p>%s</p>", err.Title, err.Message)

	io.WriteString(writer, html)
}
