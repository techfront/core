package src

import (
	"github.com/fragmenta/server"
	"github.com/techfront/core/src/kernel/assets"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/displayerror"
	"net/http"
	"os"
	"path"
	"time"
)

var appAssets *assets.Collection

/**
* Инициализация и конфигурирование Assets Pipeline.
 */
func setupAsset(server *server.Server) {
	defer server.Timef("#info Finished loading assets in %s", time.Now())

	assetsCompiled := server.ConfigBool("assets_compiled")
	appAssets = assets.New(assetsCompiled)

	err := appAssets.Load()
	if err != nil {
		server.Logf("#info Compiling assets")
		err := appAssets.Compile("tmp/src-bundle", "public")
		if err != nil {
			server.Fatalf("#error compiling assets %s", err)
		}
	}
}

/**
* Хандлер для статических файлов (js, css, etc.)
* Он расположен здесь, так как есть доступ к коллекции assets, что очень удобно.
* Это исключение. Все остальные базовые хандлеры можно найти в дирректории handler.
 */
func handlerAsset(context router.Context) error {
	var localPath string

	p := path.Clean(context.Path())
	f := appAssets.File(path.Base(p))
	if f != nil {
		localPath = "./" + f.LocalPath()

		http.ServeFile(context.Writer(), context.Request(), localPath)

		return nil
	}

	localPath = "./public" + p

	if _, err := os.Stat(localPath); err != nil {
		if os.IsNotExist(err) {
			return router.NotFoundError(err, displayerror.PageNotFound...)
		}

		return router.NotAuthorizedError(err, displayerror.AccessDenied...)
	}

	http.ServeFile(context.Writer(), context.Request(), localPath)

	return nil
}
