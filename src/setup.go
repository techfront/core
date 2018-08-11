package src

import (
	"github.com/fragmenta/server"
	"github.com/fragmenta/server/log"
	"github.com/techfront/core/src/kernel/schedule"
	"github.com/techfront/core/src/kernel/router"
	"runtime"
)

/**
 * Точка сборки всех конфигураций.
 */
func Setup(server *server.Server) {

	runtime.GOMAXPROCS(runtime.NumCPU())
	server.Logger = log.New(server.Config("log"), server.Production())
	schedule.Setup(server.Logger, server)

	setupAsset(server)
	setupView(server)
	setupService(server)
	setupDatabase(server)
	setupLib(server)

	/**
	 * Инициализация роутинга.
	 */
	r := router.New(server.Logger, server)
	setupRoute(r)
	setupMiddleware(r)
}
