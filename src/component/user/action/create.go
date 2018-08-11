package useraction

import (
	"strings"

	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"
	"github.com/techfront/core/src/lib/status"

	"github.com/techfront/core/src/component/user"
)

/**
* HandleCreateShow handles GET /signup
 */
func HandleCreateShow(context router.Context) error {
	currentUser := authorise.CurrentUser(context)

	if !currentUser.Anon() && !currentUser.IsUnverified()  {
		return router.NotAuthorizedError(nil, displayerror.AccessDenied...)
	}

	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	token := params.Get("token")

	if len(token) > 0 {
		if err := validateActivationToken(token); err != nil {
			return err
		} else {
			if err := activateUser(context, token); err != nil {
				return err
			}
		}
	}

	// Отображение шаблона
	v := view.New("component/user/template/create")

	switch params.Get("message") {
	case "error--user_exist_by_nickname":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "К сожалению, пользователь с таким же никнеймом существует, пожалуйста, выберите другой."
	case "error--user_exist_by_email":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "К сожалению, пользователь с таким же Email уже существует."
	case "error--incorrect_nickname":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Ваш никнейм слишком короткий. Пожалуйста, выберите никнейм длиной более 2 символов."
	case "error--prohibited_nickname":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Ваш никнейм содержит запрещённые слова или символы. Пожалуйста, выберите другой."
	case "error--incorrect_email":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Пожалуйста, проверьте правильность Email."
	case "error--incorrect_password":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Ваш пароль должен иметь длину не менее 8 символов. Пожалуйста, выберите другой."
	}

	userEntity := user.New()

	v.Vars["user"] = userEntity
	v.Vars["meta_title"] = "Регистрация / Техфронт"

	return v.Render(context)
}

/**
* HandleCreate handles POST /signup from the register page
 */
func HandleCreate(context router.Context) error {

	// Проверка прав
	err := authorise.AuthenticityToken(context)
	if err != nil {
		return router.NotAuthorizedError(err, displayerror.AccessDenied...)
	}

	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Проверка Email
	email := params.Get("user_email")
	if len(email) < 3 || !strings.Contains(email, "@") {
		return router.Redirect(context, "/signup?message=error--incorrect_email")
	}

	if count, err := user.Query().Where("user_email=?", email).Count(); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	} else if count > 0 {
		return router.Redirect(context, "/signup?message=error--user_exist_by_email")
	}

	// Проверка Name
	name := params.Get("user_name")
	if err := user.CheckName(name); err != nil {
		return router.Redirect(context, "/signup?message=error--prohibited_nickname")
	}
	if len(name) < 2 {
		return router.Redirect(context, "/signup?message=error--incorrect_nickname")
	}

	if count, err := user.Query().Where("user_name=?", name).Count(); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	} else if count > 0 {
		return router.Redirect(context, "/signup?message=error--user_exist_by_nickname")
	}

	// Проверка Password
	password := params.Get("user_password")
	if len(password) < 8 {
		return router.Redirect(context, "/signup?message=error--incorrect_password")
	}

	// Установка нового аватара
	params.Set("user_avatar", user.GetAvatar(email))

	// Определение базовых параметров
	params.SetInt("user_status", status.Unverified)
	params.SetInt("user_role", user.RoleReader)
	params.SetInt("user_score", 1.000)

	// Создание пользователя и очистка параметров
	id, err := user.Create(params.Clean(user.AllowedParamsAdmin()))
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение созданного пользователя по ID
	userEntity, err := user.Find(id)
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Отправка Email
	err = userEntity.SendEmailActivation()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Открытие новой сессии
	err = loginUser(context, userEntity)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return router.Redirect(context, "/")
}

/**
* validateActivationToken проверяет токен.
*
* @param token string
 */
func validateActivationToken(token string) error {
	if len(token) == 0 {
		return router.NotFoundError(nil, displayerror.PageNotFound...)
	}

	q := user.Where("user_create_token=?", token)

	userEntity, err := user.First(q)
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	if userEntity.CreateToken != token {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	return nil
}

/**
* Функция activateUser подтверждает регистрацию.
*
* @param token string
 */
func activateUser(context router.Context,token string) error {
	if len(token) == 0 {
		return router.NotFoundError(nil, displayerror.PageNotFound...)
	}

	q := user.Where("user_create_token=?", token)
	userEntity, err := user.First(q)
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	if err := userEntity.Update(map[string]string{"user_status": "100"}); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return router.Redirect(context, "/")
}
