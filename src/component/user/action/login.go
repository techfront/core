package useraction

import (
	"fmt"

	"github.com/techfront/core/src/kernel/auth"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/component/user"

	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleLoginShow shows the page at /users/login
 */
func HandleLoginShow(context router.Context) error {

	// Проверка прав
	if !authorise.CurrentUser(context).Anon() {
		return router.NotAuthorizedError(nil, displayerror.AccessDenied...)
	}

	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Отображение шаблона
	v := view.New("component/user/template/login")

	switch params.Get("message") {
	case "error--user_not_exist":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "К сожалению, пользователь не существует."
	case "error--incorrect_password":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "К сожалению, неверный пароль, повторите попытку."

	case "error--limit-exceeded":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "К сожалению, вы исчерпали лимит обращений, попробуйте снова через 1 час."
	}

	v.Vars["meta_title"] = "Войти / Техфронт"

	return v.Render(context)
}

/**
* HandleLogin handles a post to /users/login
 */
func HandleLogin(context router.Context) error {

	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Поиск пользователя по Email или Name
	query := user.Where("user_email=?", params.Get("user_email"))
	userEntity, err := user.First(query)
	if err != nil {
		query = user.Where("user_name=?", params.Get("user_email"))
		userEntity, err = user.First(query)
	}
	if err != nil {
		return router.Redirect(context, "/login?message=error--user_not_exist")
	}

	// Проверка лимита
	limitValid, err := userEntity.CheckLoginLimit()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}
	if !limitValid {
		return router.Redirect(context, "/login?message=error--limit-exceeded")
	}

	// Проверка пароля
	err = auth.CheckPassword(params.Get("user_password"), userEntity.EncryptedPassword)
	if err != nil {
		return router.Redirect(context, "/login?message=error--incorrect_password")
	}

	// Открытие новой сессии
	err = loginUser(context, userEntity)
	if err != nil {
		return router.InternalError(nil, displayerror.UnknownError...)
	}

	// Редирект на главную
	return router.Redirect(context, "/")

}

/**
* Функция declensionRemaining выбирает правильное склонение для колличества попыток.
 */
func declensionRemaining(r int64) string {
	switch {
	case (r%10 == 1) && (r%100 != 11):
		return "попытка"
	case (r%10 >= 2) && (r%10 <= 4) && (r%100 < 10 || r%100 >= 20):
		return "попытки"
	}

	return "попыток"
}

/**
* Создает новую сессию для пользователя
 */
func loginUser(context router.Context, userEntity *user.User) error {
	session, err := auth.Session(context.Writer(), context.Request())
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	session.Set(auth.SessionUserKey, fmt.Sprintf("%d", userEntity.Id))
	session.Save(context.Writer())
	return nil
}
