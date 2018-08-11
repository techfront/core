package useraction

import (
	"github.com/techfront/core/src/kernel/auth"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/response"
	"github.com/techfront/core/src/component/user"
)

func HandleAuthoriseEndpoint(context router.Context) error {

	w := context.Writer()
	r := context.Request()

	username, password, ok := r.BasicAuth()
	if !ok {
		return response.Send(w, 403, response.BadRequest("Ошибка 500. Неверные заголовки запроса."))
	}

	query := user.Where("user_name=?", username)
	userEntity, err := user.First(query)
	if err != nil {
		return response.Send(w, 401, response.Unauthorized("Ошибка 401. Пользователь не найден."))
	}

	limitValid, err := userEntity.CheckLoginLimit()
	if err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}
	if !limitValid {
		return response.Send(w, 429, response.TooManyRequests("Ошибка 429. Превышен лимит."))
	}

	if err := auth.CheckPassword(password, userEntity.EncryptedPassword); err != nil {
		return response.Send(w, 401, response.Unauthorized("Ошибка 401. Некорректный пароль."))
	}

	if err := loginUser(context, userEntity); err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}

	return response.Send(w, 200, &response.Response{
		Data: nil,
		Meta: map[string]interface{}{
			"_notice": "Авторизация успешна.",
		},
	})

}