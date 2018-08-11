package topicaction

import (
	"fmt"
	"strings"

	"github.com/techfront/core/src/kernel/validate"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"
	"github.com/techfront/core/src/lib/embedly"

	"github.com/techfront/core/src/component/topic"
)

/**
* HandleCreateShow serves the create form via GET for topics
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
	v := view.New("component/topic/template/create")

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
	case "error--topic_dublicate":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Похожий топик уже был создан ранее. Воспользуйтесь поиском, что бы найти."
	case "error--topic_invalid_create_timeout":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Создать новый топик можно через 1 минуту. Подождите и попробуйте снова."
	}

	topicEntity := topic.New()

	v.Vars["topic"] = topicEntity
	v.Vars["meta_title"] = "Cоздать новый топик / Техфронт"

	return v.Render(context)
}

/**
* HandleCreate handles the POST of the create form for topics
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
	limitValid, err := userEntity.CheckCreateTopicLimit()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}
	if !limitValid {
		return router.Redirect(context, "/submit/topic?message=error--topic_invalid_create_timeout")
	}

	// Проверка названия
	name := params.Get("topic_name")
	if err := validate.Length(name, 10, 1000); err != nil {
		return router.Redirect(context, "/submit/topic?message=error--topic_incorrect_name_length")
	}

	// Обработка URL
	url := params.Get("topic_url")
	if len(url) > 0 {
		if err := validate.Length(url, 5, 1000); err != nil {
			return router.Redirect(context, "/submit/topic?message=error--topic_incorrect_url_length")
		}

		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			return router.Redirect(context, "/submit/topic?message=error--topic_incorrect_url_prefix")
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

	// Загрузка превью
	if thumb, err := embedly.UploadEmbedlyThumbnail(url); err != nil {
		params.Remove("topic_thumbnail")
	} else {
		params.Set("topic_thumbnail", thumb)
	}

	// Проверка топика на существование
	query := topic.Where("topic_url=?", url)
	duplicates, err := topic.FindAll(query)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	if len(duplicates) > 0 {
		topicEntity := duplicates[0]

		if topicHasUserVote(topicEntity, userEntity) {
			return router.Redirect(context, "/submit/topic?message=error--topic_dublicate")
		}

		topicUpvote(topicEntity, userEntity, ip)

		return router.Redirect(context, topicEntity.URLShow())
	}

	// Очистка параметров для текущей роли
	accepted := topic.AllowedParams()
	if userEntity.Admin() {
		accepted = topic.AllowedParamsAdmin()
	}
	cleanedParams := params.Clean(accepted)

	// Добавление базовых параметров
	cleanedParams["topic_count_upvote"] = "1"
	cleanedParams["topic_id_user"] = fmt.Sprintf("%d", userEntity.Id)
	cleanedParams["topic_status"] = "14"

	if userEntity.CanPublish() {
		cleanedParams["topic_status"] = "100"
	}

	// Создание топика
	id, err := topic.Create(cleanedParams)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение созданного топика по ID
	topicEntity, err := topic.Find(id)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Обновление колличества топиков пользователя
	userParams := map[string]string{"user_count_topic": fmt.Sprintf("%d", userEntity.TopicCount+1)}
	err = userEntity.Update(userParams)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Ставим первый голос
	err = topicRecordUpvote(topicEntity, userEntity, ip)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Редирект на главную, если пользователь опубликовал топик
	if userEntity.CanPublish() {
		return router.Redirect(context, "/")
	}

	// Редирект к топику, если топик помещён на проверку
	return router.Redirect(context, topicEntity.URLShow()+"/?message=info--topic_in_the_queue")
}
