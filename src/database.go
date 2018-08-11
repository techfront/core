package src

import (
	"github.com/fragmenta/query"
	"github.com/fragmenta/server"
	"time"
)

/**
* Инициализация и конфигурирование базы данных.
 */
func setupDatabase(server *server.Server) {
	defer server.Timef("#info Finished opening in %s database %s for user %s", time.Now(), server.Config("db"), server.Config("db_user"))

	config := server.Configuration()

	options := map[string]string{
		"adapter":  config["db_adapter"],
		"user":     config["db_user"],
		"password": config["db_pass"],
		"db":       config["db"],
	}

	err := query.OpenDatabase(options)

	if err != nil {
		server.Fatalf("Error reading database %s", err)
	}
}
