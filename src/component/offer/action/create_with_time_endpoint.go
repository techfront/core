package offeraction

import (
	"fmt"
	"time"
	"github.com/techfront/core/src/kernel/schedule"
	"github.com/techfront/core/src/kernel/validate"

	"github.com/techfront/core/src/component/offer"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/authorise"

	"github.com/techfront/core/src/kernel/response"
)

/**
 * HandleCreateWithTimeEndpoint добавляет оффер.
 *
 */
func HandleCreateWithTimeEndpoint(context router.Context) error {
	w := context.Writer()

	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}

	timer := params.Get("offer_publication_timer")
	timerDuration, err := time.ParseDuration(timer)
	if err != nil {
		return response.Send(w, 403, response.Forbidden("Ошибка 403. Некорректный параметр offer_publication_timer."))
	}
	publicationTime := time.Now().Add(timerDuration)

	// Получение пользователя
	userEntity := authorise.CurrentUser(context)
	ip := context.ClientIP()

	// Проверка прав
	if !userEntity.Admin() {
		return response.Send(w, 401, response.Unauthorized("Ошибка 401. К сожалению, у вас плохая карма."))
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
	cleanedParams["offer_status"] = "14"

	if userEntity.CanPublish() {
		cleanedParams["offer_status"] = "100"
	}

	// Создание задания
	scheduleContext := schedule.NewContext()
	scheduleContext.Set("cleaned_params", cleanedParams)
	scheduleContext.Set("ip", ip)

	createFunc := func(context schedule.Context) {
		id, err := offer.Create(context.Get("cleaned_params").(map[string]string))
		if err == nil {
			offerEntity, err := offer.Find(id)
			if err == nil {
				userParams := map[string]string{"user_count_offer": fmt.Sprintf("%d", userEntity.OfferCount + 1)}
				userEntity.Update(userParams)
				offerRecordUpvote(offerEntity, userEntity, context.Get("ip").(string)) 
			}
		}
    	}

    	// Добавление в очередь
    	schedule.At(createFunc, scheduleContext, publicationTime, 0)

	return response.Send(w, 200, &response.Response{
		Data: nil,
		Meta: map[string]interface{}{
			"_notice": "Оффер добавлен в очередь и скоро будет опубликован.",
		},
	})
}
