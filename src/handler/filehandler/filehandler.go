package filehandler

import (
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/displayerror"
	"net/http"
	"os"
	"path"
)

func Handler(context router.Context) error {
	localPath := "./public" + path.Clean(context.Path())

	if _, err := os.Stat(localPath); err != nil {
		if os.IsNotExist(err) {
			return router.NotFoundError(err, displayerror.PageNotFound...)
		}

		return router.NotAuthorizedError(err, displayerror.AccessDenied...)
	}

	http.ServeFile(context.Writer(), context.Request(), localPath)

	return nil
}
