package topicaction

import (
	"strings"

	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/validate"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"
	"github.com/techfront/core/src/lib/upload"

	"github.com/techfront/core/src/component/topic"
)

/**
* HandleUpdateShow renders the form to update a topic
 */
func HandleUpdateShow(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение топика
	topicEntity, err := topic.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Если топик отклонён, то обновить нельзя
	if topicEntity.IsRejected() && !authorise.CurrentUser(context).Admin() {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Проверка прав
	err = authorise.Resource(context, topicEntity)
	if err != nil {
		return router.NotAuthorizedError(err, displayerror.AccessDenied...)
	}

	// Отображение шаблона
	v := view.New(
		"component/topic/template/update",
		"component/topic/template/form",
	)

	switch params.Get("message") {
	case "error--topic_incorrect_name_length":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Заголовок должен содержать не менее 10, но не более 150 символов."
	case "error--topic_incorrect_url_length":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Ссылка должна содержать более пяти символов."
	case "error--topic_incorrect_url_prefix":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Ссылка должна начинаться с https:// или http:// ."
	}

	v.Vars["topic"] = topicEntity
	v.Vars["meta_title"] = "Обновление топика / Техфронт"

	return v.Render(context)
}

/**
* HandleUpdate handles the POST of the form to update a topic
 */
func HandleUpdate(context router.Context) error {

	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение топика
	topicEntity, err := topic.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Проверка прав
	err = authorise.ResourceAndAuthenticity(context, topicEntity)
	if err != nil {
		return router.NotAuthorizedError(err, displayerror.AccessDenied...)
	}

	// Проверка заголовка
	name := params.Get("topic_name")
	if err := validate.Length(name, 10, 1000); err != nil {
		return router.Redirect(context, "/topics/"+params.Get("id")+"/update?message=error--topic_incorrect_name_length")
	}

	// Обработка URL-ов
	url := params.Get("topic_url")
	if len(url) > 0 {
		if err := validate.Length(url, 5, 1000); err != nil {
			return router.Redirect(context, "/topics/"+params.Get("id")+"/update?message=error--topic_incorrect_url_length")
		}

		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			return router.Redirect(context, "/topics/"+params.Get("id")+"/update?message=error--topic_incorrect_url_prefix")
		}
	}

	if strings.HasSuffix(url, "/") {
		url = strings.Trim(url, "/")
	}

	if strings.Contains(url, "?utm_") {
		url = strings.Split(url, "?utm_")[0]
	}

	if strings.Contains(url, "#") {
		url = strings.Split(url, "#")[0]
	}

	params.Set("topic_url", url)

	// Загружаем Thumbnail, если в поле находится ссылка
	thumbnail := params.Get("topic_thumbnail")
	if (strings.HasPrefix(thumbnail, "http://") || strings.HasPrefix(thumbnail, "https://")) && authorise.CurrentUser(context).Admin() {
		if newThumbnail, err := upload.UploadFromUrl(thumbnail, "thumbnail"); err != nil {
			params.Remove("topic_thumbnail")
		} else {
			params.Set("topic_thumbnail", newThumbnail)
		}
	}

	// Очистка параметров для текущей роли
	accepted := topic.AllowedUpdateParams()
	if authorise.CurrentUser(context).Admin() {
		accepted = topic.AllowedParamsAdmin()
	}
	cleanedParams := params.Clean(accepted)

	err = topicEntity.Update(cleanedParams)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Редирект к топику
	return router.Redirect(context, topicEntity.URLShow())
}
