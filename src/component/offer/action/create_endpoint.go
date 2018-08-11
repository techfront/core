package offeraction

import (
	"fmt"
	"github.com/techfront/core/src/kernel/validate"

	"github.com/techfront/core/src/component/offer"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/authorise"

	"github.com/techfront/core/src/kernel/response"
)

/**
 * HandleCreateEndpoint добавляет оффер.
 *
 */
func HandleCreateEndpoint(context router.Context) error {

	w := context.Writer()

	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}

	// Получение пользователя
	userEntity := authorise.CurrentUser(context)
	ip := context.ClientIP()

	// Проверка прав
	if !userEntity.CanSubmit() {
		return response.Send(w, 401, response.Unauthorized("Ошибка 401. К сожалению, у вас плохая карма."))
	}

	// Проверка лимита
	limitValid, err := userEntity.CheckCreateOfferLimit()
	if err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}
	if !limitValid {
		return response.Send(w, 429, response.TooManyRequests("Ошибка 429. Превышен лимит."))
	}

	// Проверка названия
	name := params.Get("offer_name")
	if err := validate.Length(name, 10, 1000); err != nil {
		return response.Send(w, 403, response.Forbidden("Ошибка 403. Некорректный параметр offer_name."))
	}

	// Очистка параметров для текущей роли
	accepted := offer.AllowedParams()
	if userEntity.Admin() {
		accepted = offer.AllowedParamsAdmin()
	}
	cleanedParams := params.Clean(accepted)

	// Добавление базовых параметров
	cleanedParams["offer_count_upvote"] = "1"
	cleanedParams["offer_id_user"] = fmt.Sprintf("%d", userEntity.Id)
	cleanedParams["offer_status"] = "100"

	// Создание оффера
	id, err := offer.Create(cleanedParams)
	if err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}

	// Получение созданного оффера по ID
	offerEntity, err := offer.Find(id)
	if err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}

	// Обновление колличества офферов пользователя
	userParams := map[string]string{"user_count_offer": fmt.Sprintf("%d", userEntity.OfferCount+1)}
	err = userEntity.Update(userParams)
	if err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}

	// Ставим первый голос
	err = offerRecordUpvote(offerEntity, userEntity, ip)
	if err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}

	var notice string

	if userEntity.CanPublish() {
		notice = "Оффер успешно создан."
	} else {
		notice = "Оффер успешно добавлен в очередь на проверку."
	}

	return response.Send(w, 200, &response.Response{
		Data: nil,
		Meta: map[string]interface{}{
			"_notice": notice,
		},
	})
}
