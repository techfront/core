package useraction

import (
	"fmt"
	"github.com/techfront/core/src/component/user"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/response"
	"github.com/techfront/core/src/lib/authorise"
)

/**
* HandleActivateEndpoint отправляет письмо активации пользователя на Email.
 */
func HandleActivateEndpoint(context router.Context) error {

	w := context.Writer()

	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}

	// Получение пользователя
	userEntity, err := user.Find(params.GetInt("id"))
	if err != nil {
		return response.Send(w, 404, response.NotFound("Ошибка 404. Ресурс не найден."))
	}

	// Проверка прав
	err = authorise.ResourceAndAuthenticity(context, userEntity)
	if err != nil {
		return response.Send(w, 401, response.Unauthorized("Ошибка 401. Необходима авторизация."))
	}

	// Проверка лимита
	limitValid, remaining, err := userEntity.CheckActivationLimit()
	if err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}
	if !limitValid {
		return response.Send(w, 429, response.TooManyRequests("Ошибка 429. Запросов слишком много, попробуйте повторить через час или обратитесь к администратору."))
	}

	// Отправка письма
	err = userEntity.SendEmailActivation()
	if err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}

	return response.Send(w, 200, &response.Response{
		Data: nil,
		Meta: map[string]interface{}{
			"_notice": fmt.Sprintf("Письмо успешно отправлено. Осталось попыток: %d/4", remaining),
		},
	})
}