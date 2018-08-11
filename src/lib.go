package src

import (
	"github.com/fragmenta/server"
	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/cache"
	"github.com/techfront/core/src/lib/embedly"
	"github.com/techfront/core/src/lib/gzip"
	"github.com/techfront/core/src/lib/limiter"
	"github.com/techfront/core/src/lib/mail"
	"github.com/techfront/core/src/lib/mailchimp"
	"github.com/techfront/core/src/lib/media"
	"github.com/techfront/core/src/lib/resizer"
	"github.com/techfront/core/src/lib/upload"
	"github.com/techfront/core/src/lib/urlshortener"
)

/**
* Инициализация и конфигурирование библиотек.
 */
func setupLib(server *server.Server) {

	config := server.Configuration()

	/**
	 * Конфигурирование авторизации.
	 */
	authorise.Setup(server)

	/**
	 * Конфигурирование Media (Размеры и форматы медиа-файлов).
	 */
	media.Setup()

	/**
	 * Конфигурирование кэширования.
	 */
	cache.Setup(config)

	/**
	 * Конфигурирование Google Shorten API.
	 */
	urlshortener.Setup(config)

	/**
	 * Конфигурирование почтового сервера.
	 */
	mail.Setup(config)

	/**
	 * Конфигурирование загрузчика изображений.
	 */
	upload.Setup(config)

	/**
	 * Конфигурирование ресайзера изображений.
	 */
	resizer.Setup(config)

	/**
	 * Конфигурирование Embedly API.
	 */
	embedly.Setup(config)

	/**
	 * Конфигурирование Gzip.
	 */
	gzip.Setup()

	/**
	 * Конфигурирование Limiter (Алгоритм базовой скорости ячеек).
	 */
	limiter.Setup(config)

	/**
	 * Конфигурирование Mailchimp API.
	 */
	mailchimp.Setup(config)
}