package useraction

import (
	"time"

	"github.com/techfront/core/src/kernel/auth"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"
	"github.com/techfront/core/src/lib/mail"

	"github.com/techfront/core/src/component/user"
)

func HandleResetShow(context router.Context) error {

	if !authorise.CurrentUser(context).Anon() {
		return router.NotAuthorizedError(nil, displayerror.AccessDenied...)
	}

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	token := params.Get("token")

	// Отображение шаблона
	v := view.New("component/user/template/reset")

	if len(token) > 0 {
		if err := validateResetToken(token); err != nil {
			return err
		} else {
			v = view.New("component/user/template/change_password")

			v.Vars["user_reset_token"] = token
		}
	}

	switch params.Get("message") {
	case "error--user_not_exist":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "К сожалению, пользователь не существует."
	case "error--incorrect_password":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Пароль должен содержать не менее 8 символов, повторите попытку."
	case "info--email_sent":
		v.Vars["message_type"] = "default"
		v.Vars["message"] = "Письмо с инструкцией отправлено, проверьте вашу почту."
	case "success--password_changed":
		v.Vars["message_type"] = "success"
		v.Vars["message"] = "Пароль успешно изменён."
	}

	v.Vars["meta_title"] = "Восстановление доступа / Техфронт"

	return v.Render(context)
}

func HandleReset(context router.Context) error {
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Изменить пароль
	if len(params.Get("user_password")) > 0 {
		return changePassword(context)
	}

	// Поиск пользователя по Email
	email := params.Get("user_email")
	userEntity, err := user.FindEmail(email)
	if err != nil {
		return router.Redirect(context, "/reset?message=error--user_not_exist")
	}

	// Создание токена
	token := auth.BytesToHex(auth.RandomToken())

	// Обновление сущности пользователя
	if err = userEntity.Update(map[string]string{"user_reset_token": token}); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Отправка Email
	recipients := []string{email}
	variables := map[string]interface{}{"reset_link": "https://techfront.org/reset?token=" + token, "user_name": userEntity.Name}
	if err := mail.Send(recipients, "Восстановление доступа", "component/user/template/mail/reset", variables, nil); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return router.Redirect(context, "/reset?message=info--email_sent")
}

func changePassword(context router.Context) error {
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	token := params.Get("user_reset_token")

	// Проверка пароля
	if len(params.Get("user_password")) < 8 {
		return router.Redirect(context, "/reset?message=error--incorrect_password")
	}

	// Проверка ключа
	if err := validateResetToken(token); err != nil {
		return err
	}

	// Поиск пользователя по Токену
	q := user.Where("user_reset_token=?", token)

	userEntity, err := user.First(q)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Обнуление токена
	params.Set("user_reset_token", "")

	// Очистка от лишних параметров для обновления
	allowedParams := params.Clean([]string{"user_password", "user_reset_token"})

	// Обновление сущности пользователя
	if err = userEntity.Update(allowedParams); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return router.Redirect(context, "/reset?message=success--password_changed")

}

/**
* validateResetToken проверяет токен.
*
* @param token string
 */
func validateResetToken(token string) error {
	// Если параметр пуст, то возвращает ошибку.
	if len(token) == 0 {
		return router.NotFoundError(nil, displayerror.PageNotFound...)
	}

	// Поиск пользователя
	q := user.Where("user_reset_token=?", token)

	userEntity, err := user.First(q)
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Проверка ключа
	if userEntity.ResetToken != token {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Проверка срока истечения
	duration := time.Since(userEntity.UpdatedAt)
	minutes := duration / time.Minute
	if minutes > 5 {
		return router.InternalError(nil, displayerror.TimeIsOverError...)
	}

	return nil
}
