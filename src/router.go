package src

import (
	"github.com/techfront/core/src/kernel/router"

	"github.com/techfront/core/src/handler/errorhandler"
	"github.com/techfront/core/src/handler/filehandler"
	// "github.com/techfront/core/src/handler/pprofhandler"
	"github.com/techfront/core/src/handler/redirecthandler"

	"github.com/techfront/core/src/middleware/ahtnmiddleware"
	"github.com/techfront/core/src/middleware/cumiddleware"
	"github.com/techfront/core/src/middleware/cuomiddleware"
	"github.com/techfront/core/src/middleware/logmiddleware"
	"github.com/techfront/core/src/middleware/ascemiddleware"
	"github.com/techfront/core/src/middleware/gzipmiddleware"
	"github.com/techfront/core/src/middleware/secmiddleware"

	"github.com/techfront/core/src/component/admin/action"
	"github.com/techfront/core/src/component/comment/action"
	"github.com/techfront/core/src/component/page/action"
	"github.com/techfront/core/src/component/topic/action"
	"github.com/techfront/core/src/component/offer/action"
	"github.com/techfront/core/src/component/user/action"
)

/**
 * Конфигурирование маршрутов.
 */
func setupRoute(r *router.Router) {

	// Обработка ошибок
	r.ErrorHandler = errorhandler.Handler

	// Обработка файлов
	r.FileHandler = filehandler.Handler

	// Блок маршрутов для статики
	r.Get("/uploads/*filepath", filehandler.Handler)
	r.Get("/images/*filepath", filehandler.Handler)
	r.Get("/assets/*filepath", handlerAsset)

	// Блок маршрутов для топиков
	r.Get("/", topicaction.HandleHome)
	r.Get("/topics", topicaction.HandleIndex)
	r.Get("/xml/topics", topicaction.HandleIndexXml)
	r.Get("/xml/top", topicaction.HandleTopXml)

	r.Get("/format/:format", topicaction.HandleFormat)
	r.Get("/xml/format/:format", topicaction.HandleFormatXml)

	r.Get("/submit/topic", topicaction.HandleCreateShow)
	r.Post("/submit/topic", topicaction.HandleCreate)

	r.Get("/topics/:id", topicaction.HandleShow)
	r.Get("/topics/:id/update", topicaction.HandleUpdateShow)
	r.Post("/topics/:id/update", topicaction.HandleUpdate)
	r.Post("/topics/:id/destroy", topicaction.HandleDestroy)
	r.Post("/topics/:id/upvote", topicaction.HandleUpvote)
	r.Post("/topics/:id/downvote", topicaction.HandleDownvote)

	r.Get("/favorited", topicaction.HandleFavorited)


	// Блок маршрутов для офферов
	/*
	r.Get("/underground", offeraction.HandleIndex)
	r.Get("/offers", topicaction.HandleIndex)
	
	r.Get("/submit/offer", offeraction.HandleCreateShow)
	r.Post("/submit/offer", offeraction.HandleCreate)

	r.Get("/offers/:id", offeraction.HandleShow)
	r.Get("/offers/:id/update", offeraction.HandleUpdateShow)
	r.Post("/offers/:id/update", offeraction.HandleUpdate)
	r.Post("/offers/:id/destroy", offeraction.HandleDestroy)
	r.Post("/offers/:id/upvote", offeraction.HandleUpvote)
	r.Post("/offers/:id/downvote", offeraction.HandleDownvote)

	r.Get("/underground/favorited", offeraction.HandleFavorited)
	*/

	// Блок маршрутов для комментариев
	r.Get("/discussions", commentaction.HandleIndex)

	r.Post("/submit/comment", commentaction.HandleCreate)

	r.Get("/comments/:id/update", commentaction.HandleUpdateShow)
	r.Post("/comments/:id/update", commentaction.HandleUpdate)
	r.Post("/comments/:id/destroy", commentaction.HandleDestroy)
	r.Post("/comments/:id/hide", commentaction.HandleHide)

	r.Post("/comments/:id/upvote", commentaction.HandleUpvote)
	r.Post("/comments/:id/downvote", commentaction.HandleDownvote)

	// Блок маршрутов для пользователей
	r.Get("/users", useraction.HandleIndex)
	r.Get("/users/:id", useraction.HandleShow)

	r.Get("/login", useraction.HandleLoginShow)
	r.Post("/login", useraction.HandleLogin)

	r.Post("/logout", useraction.HandleLogout)

	r.Get("/reset", useraction.HandleResetShow)
	r.Post("/reset", useraction.HandleReset)

	r.Get("/signup", useraction.HandleCreateShow)
	r.Post("/signup", useraction.HandleCreate)

	r.Get("/unsubscribe", useraction.HandleUnsubscribeShow)
	r.Post("/unsubscribe", useraction.HandleUnsubscribe)

	r.Get("/users/:id/update", useraction.HandleUpdateShow)
	r.Post("/users/:id/update", useraction.HandleUpdate)
	r.Post("/users/:id/destroy", useraction.HandleDestroy)

	// Блок маршрутов статических страниц
	r.Get("/about", pageaction.HandleAbout)
	r.Get("/ways-to-give", pageaction.HandleWaysToGive)
	r.Get("/spasibo", pageaction.HandleSpasibo)

	// Блок маршрутов Администратора
	r.Get("/admin/topics", adminaction.HandleTopic)

	// Конечные точки API
	r.Get("/api/v3/get/topic/feed", topicaction.HandleFeedEndpoint)
	r.Post("/api/v3/post/topic/create", topicaction.HandleCreateWithTimeEndpoint)

	// r.Post("/api/v3/post/user/auth", useraction.HandleAuthoriseEndpoint)
	r.Post("/api/v3/post/user/activate", useraction.HandleActivateEndpoint)
	r.Post("/api/v3/post/user/subscribe", useraction.HandleSubscribeEndpoint)
	r.Post("/api/v3/post/user/favorite", useraction.HandleCreateFavoriteEndpoint)
	r.Post("/api/v3/post/user/unfavorite", useraction.HandleDestroyFavoriteEndpoint)

	// Debug
	// r.Get("/debug/pprof/*pprof", pprofhandler.Handler)

	// Redirect
	r.Get("/away", redirecthandler.Handler)
}

/**
 * Настройка посредников.
 */
func setupMiddleware(r *router.Router) {

	// Логирует запрос.
	r.AddMiddleware(logmiddleware.Middleware)

	// Кеширует Get-запрос для анонимуса
	r.AddMiddleware(ascemiddleware.Middleware)

	// Добавляет сущность текущего пользователя в контекст
	r.AddMiddleware(cumiddleware.Middleware)

	// Добавляет токен в контекст для аутентификации
	r.AddMiddleware(ahtnmiddleware.Middleware)

	// Проверяет статус визита пользователя
	r.AddMiddleware(cuomiddleware.Middleware)

	// Добавляет заголовки безопасности
	r.AddMiddleware(secmiddleware.Middleware)

	// Сжимает ресурсы посредством Gzip
	r.AddMiddleware(gzipmiddleware.Middleware)
}
