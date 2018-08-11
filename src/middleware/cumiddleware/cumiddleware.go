package cumiddleware

import (
	"fmt"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/authorise"
)

/**
* Посредник получает и добавляет в контекст сущность текущего пользователя.
*
* @param next router.Handler следующий обработчик.
 */
func Middleware(next router.Handler) router.Handler {
	return router.HandlerFunc(func(context router.Context) error {

		if should(context) {
			userEntity := authorise.CurrentUser(context)
			userIDString := fmt.Sprintf("%d", userEntity.Id)

			// Передача объекта в контекст.
			context.Set("current_user", userEntity)

			// Передача в контекст сообщения, если пользователь не верефицирован.
			if userEntity.IsUnverified() {
				uid := "user-activate"
				msg := `<img class="content__message-thumb" src="/images/icons/32/open-letter.png" srcset="/images/icons/64/open-letter.png 2x" width="32" height="32"/>
					<p class="content__message-info">Для активации аккаунта на указанный вами email было отправлено сообщение! Не получили? Проверьте папку "Спам" либо <a class="js-message--` + uid + `-send" data-id-user="` + userIDString + `" href="#">повторите отправку</a>.</p>
					<a href="#" class="js-message-close content__message-close" data-message=".js-message--` + uid + `">×</a>`

				context.SetMessage(uid, msg, router.MessageWarning)
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
	r := context.Request()
	if r.Method != "GET" {
		return false
	}

	ext := context.Ext()
	if len(ext) > 0 {
		return false
	}

	return true
}