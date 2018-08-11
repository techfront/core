package offeraction

import (

	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/validate"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"

	"github.com/techfront/core/src/component/offer"
)

/**
* HandleUpdateShow renders the form to update a offer
 */
func HandleUpdateShow(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение оффера
	offerEntity, err := offer.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Если оффер отклонён, то обновить нельзя
	if offerEntity.IsRejected() && !authorise.CurrentUser(context).Admin() {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Проверка прав
	err = authorise.Resource(context, offerEntity)
	if err != nil {
		return router.NotAuthorizedError(err, displayerror.AccessDenied...)
	}

	// Отображение шаблона
	v := view.New(
		"component/offer/template/update",
		"component/offer/template/form",
	)

	switch params.Get("message") {
	case "error--offer_incorrect_name_length":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Заголовок должен содержать не менее 10, но не более 150 символов."
	}

	v.Vars["offer"] = offerEntity
	v.Vars["meta_title"] = "Обновление оффера / Техфронт"

	return v.Render(context)
}

/**
* HandleUpdate handles the POST of the form to update a offer
 */
func HandleUpdate(context router.Context) error {

	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение оффера
	offerEntity, err := offer.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Проверка прав
	err = authorise.ResourceAndAuthenticity(context, offerEntity)
	if err != nil {
		return router.NotAuthorizedError(err, displayerror.AccessDenied...)
	}

	// Проверка заголовка
	name := params.Get("offer_name")
	if err := validate.Length(name, 10, 1000); err != nil {
		return router.Redirect(context, "/offers/"+params.Get("id")+"/update?message=error--offer_incorrect_name_length")
	}

	// Очистка параметров для текущей роли
	accepted := offer.AllowedUpdateParams()
	if authorise.CurrentUser(context).Admin() {
		accepted = offer.AllowedParamsAdmin()
	}
	cleanedParams := params.Clean(accepted)

	err = offerEntity.Update(cleanedParams)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Редирект к офферу
	return router.Redirect(context, offerEntity.URLShow())
}
