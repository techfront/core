package useraction

import (
	"strings"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"
	"github.com/techfront/core/src/component/user"
	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleUpdateShow serves a get request at /users/1/update (show form to update)
 */
func HandleUpdateShow(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение пользователя
	userEntity, err := user.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Проверка прав
	err = authorise.Resource(context, userEntity)
	if err != nil {
		return router.NotAuthorizedError(err, displayerror.AccessDenied...)
	}

	// Отображение шаблона
	v := view.New(
		"component/user/template/update",
		"component/user/template/form",
	)

	switch params.Get("message") {
	case "error--user_exist_by_nickname":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "К сожалению, пользователь с таким никнеймом уже существует. Пожалуйста, выберите другой."
	case "error--user_exist_by_email":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "К сожалению, пользователь с таким Email уже существует."
	case "error--incorrect_nickname":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Ваш никнейм слишком короткий. Пожалуйста, выберите никнейм длиной более 2 символов."
	case "error--prohibited_nickname":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Ваш никнейм содержит запрещённые слова или символы. Пожалуйста, выберите другой."
	case "error--incorrect_email":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Ошибка при обработке E-mail. Пожалуйста, проверьте корректность."
	case "error--incorrect_password":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Ваш пароль должен иметь длину не менее 8 символов. Пожалуйста, выберите другой."
	case "error--incorrect_contacts":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Ошибка при обработке контактов. Пожалуйста, проверьте контакты на корректность."
	}

	v.Vars["user"] = userEntity
	v.Vars["meta_title"] = "Настройки профиля / Техфронт"

	return v.Render(context)
}

/**
* HandleUpdate or PUT /users/1/update
 */
func HandleUpdate(context router.Context) error {
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение пользователя
	userId := params.Get("id")
	userEntity, err := user.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Проверка прав
	err = authorise.ResourceAndAuthenticity(context, userEntity)
	if err != nil {
		return router.NotAuthorizedError(err, displayerror.AccessDenied...)
	}

	// Проверка Email
	email := params.Get("user_email")
	if email != userEntity.Email && len(email) != 0 {
		if len(email) < 3 || !strings.Contains(email, "@") {
			return router.Redirect(context, "/users/" + userId + "/update?message=error--incorrect_email")
		}
		if count, err := user.Query().Where("user_email=?", email).Count(); err != nil {
			return router.InternalError(err, displayerror.UnknownError...)
		} else if count > 0 {
			return router.Redirect(context, "/users/" + userId + "/update?message=error--user_exist_by_email")
		}
	} else {
		params.Remove("user_email")
	}

	// Установка нового аватара
	if email != userEntity.Email && len(email) != 0 {
		params.Set("user_avatar", user.GetAvatar(email))
	} else {
		params.Remove("user_avatar")
	}

	// Проверка имени
	name := params.Get("user_name")
	if name != userEntity.Name && len(name) != 0 {
		if err := user.CheckName(name); err != nil {
			return router.Redirect(context, "/users/" + userId + "/update?message=error--prohibited_nickname")
		}
		if len(name) < 2 {
			return router.Redirect(context, "/users/" + userId + "/update?message=error--incorrect_nickname")
		}

		if count, err := user.Query().Where("user_name=?", name).Count(); err != nil {
			return router.InternalError(err, displayerror.UnknownError...)
		} else if count > 0 {
			return router.Redirect(context, "/users/" + userId + "/update?message=error--user_exist_by_nickname")
		}
	} else {
		params.Remove("user_name")
	}

	// Проверка пароля
	password := params.Get("user_password")
	if len(password) == 0 {
		params.Remove("user_password")
	} else if len(password) < 8 {
		return router.Redirect(context, "/users/" + userId + "/update?message=error--incorrect_password")
	}

	// Обработка параметров контактов
	var contactParams []map[string]string

	for _, p := range user.AllowedContactParams() {
		contactValues := params.GetAll(p)
		contactName := p[len(user.CONTACT_NAME_PREFIX):]

		for _, v := range contactValues {
			if len(v) > 0 {
				err := user.ValidateContact(v)
				if err != nil {
					return router.Redirect(context, "/users/" + userId + "/update?message=error--incorrect_contacts")
				}

				contactParams = append(contactParams, map[string]string{
					"user_contact_id_user": userId,
					"user_contact_name": contactName,
					"user_contact_value": v,
				})
			}
		}
	}

	// Удаляем все старые контакты пользователя
	contactList := userEntity.FindAllContacts()
	for _, c := range contactList {
		err := c.DestroyContact()
		if err != nil {
			return router.InternalError(err, displayerror.UnknownError...)
		}
	}

	// Записываем новые контакты
	for _, p := range contactParams {
		_, err := user.CreateContact(p)
		if err != nil {
			return router.InternalError(err, displayerror.UnknownError...)
		}
	}

	// Очистка параметров для текущей роли
	accepted := user.AllowedUpdateParams()
	if authorise.CurrentUser(context).Admin() {
		accepted = user.AllowedParamsAdmin()
	}

	allowedParams := params.Clean(accepted)
	err = userEntity.Update(allowedParams)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Редирект на страницу пользователя
	return router.Redirect(context, userEntity.URLShow())
}
