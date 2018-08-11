package offeraction

import (
	"fmt"

	"github.com/techfront/core/src/kernel/validate"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"

	"github.com/techfront/core/src/component/offer"
)

/**
* HandleCreateShow serves the create form via GET for offers
 */
func HandleCreateShow(context router.Context) error {

	// Проверка прав
	err := authorise.Path(context)
	if err != nil {
		router.Redirect(context, "/signup")
	}

	// Получение прав
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	if !authorise.CurrentUser(context).CanSubmit() {
		return router.NotAuthorizedError(nil, displayerror.AccessDeniedKarma...)
	}

	// Отображение шаблона
	v := view.New("component/offer/template/create")

	switch params.Get("message") {
	case "error--offer_incorrect_name_length":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Заголовок должен содержать не менее 10, но не более 300 символов."
	case "error--offer_invalid_create_timeout":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Создать новый оффер можно через 1 минуту. Подождите и попробуйте снова."
	}

	offerEntity := offer.New()

	v.Vars["offer"] = offerEntity
	v.Vars["meta_title"] = "Cоздать новый оффер / Техфронт"

	return v.Render(context)
}

/**
* HandleCreate handles the POST of the create form for offers
 */
func HandleCreate(context router.Context) error {

	// Проверка прав
	err := authorise.AuthenticityToken(context)
	if err != nil {
		return router.NotAuthorizedError(err, displayerror.AccessDenied...)
	}

	if !authorise.CurrentUser(context).CanSubmit() {
		return router.NotAuthorizedError(nil, displayerror.AccessDeniedKarma...)
	}

	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение пользователя
	userEntity := authorise.CurrentUser(context)
	ip := context.ClientIP()

	// Проверка лимита
	limitValid, err := userEntity.CheckCreateOfferLimit()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}
	if !limitValid {
		return router.Redirect(context, "/submit/offer?message=error--offer_invalid_create_timeout")
	}

	// Проверка названия
	name := params.Get("offer_name")
	if err := validate.Length(name, 10, 1000); err != nil {
		return router.Redirect(context, "/submit/offer?message=error--offer_incorrect_name_length")
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
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение созданного оффера по ID
	offerEntity, err := offer.Find(id)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Обновление колличества офферов пользователя
	userParams := map[string]string{"user_count_offer": fmt.Sprintf("%d", userEntity.OfferCount+1)}
	err = userEntity.Update(userParams)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Ставим первый голос
	err = offerRecordUpvote(offerEntity, userEntity, ip)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Редирект на главную, если пользователь опубликовал оффер
	if userEntity.CanPublish() {
		return router.Redirect(context, "/")
	}

	// Редирект к офферу, если оффер помещён на проверку
	return router.Redirect(context, offerEntity.URLShow()+"/?message=info--offer_in_the_queue")
}
