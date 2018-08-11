package topicaction

import (
	"log"
	"fmt"
	"time"
	"strings"
	"github.com/techfront/core/src/kernel/schedule"
	"github.com/techfront/core/src/kernel/validate"
	"github.com/techfront/core/src/lib/upload"

	"github.com/techfront/core/src/component/topic"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/authorise"

	"github.com/techfront/core/src/kernel/response"
)

/**
 * HandleCreateWithTimeEndpoint добавляет топик.
 *
 */
func HandleCreateWithTimeEndpoint(context router.Context) error {

	w := context.Writer()

	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}

	timer := params.Get("topic_publication_timer")
	timerDuration, err := time.ParseDuration(timer)
	if err != nil {
		return response.Send(w, 403, response.Forbidden("Ошибка 403. Некорректный параметр topic_publication_timer."))
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
	name := params.Get("topic_name")
	log.Print(name)
	if err := validate.Length(name, 10, 1000); err != nil {
		return response.Send(w, 403, response.Forbidden("Ошибка 403. Некорректный параметр topic_name."))
	}

	// Обработка URL
	url := params.Get("topic_url")
	if len(url) > 0 {
		if err := validate.Length(url, 5, 1000); err != nil {
			return response.Send(w, 403, response.Forbidden("Ошибка 403. Некорректный параметр topic_url."))
		}

		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			return response.Send(w, 403, response.Forbidden("Ошибка 403. Некорректный параметр topic_url."))
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

	// Проверка топика на существование
	query := topic.Where("topic_url=?", url)
	duplicates, err := topic.FindAll(query)
	if err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}

	if len(duplicates) > 0 {
		return response.Send(w, 403, response.Forbidden("Ошибка 403. Топик был создан ранее."))
	}

	// Загрузка превью
	thumbnail := params.Get("topic_thumbnail")
	if len(thumbnail) > 0 {
		path, err := upload.UploadFromUrl(thumbnail, "thumbnail")
		if err != nil {
			log.Print(err)
			return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
		}
		params.Set("topic_thumbnail", path)
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

	// Создание задания
	scheduleContext := schedule.NewContext()
	scheduleContext.Set("cleaned_params", cleanedParams)
	scheduleContext.Set("ip", ip)

	createFunc := func(context schedule.Context) {
		id, err := topic.Create(context.Get("cleaned_params").(map[string]string))
		if err == nil {
			topicEntity, err := topic.Find(id)
			if err == nil {
				userParams := map[string]string{"user_count_topic": fmt.Sprintf("%d", userEntity.TopicCount + 1)}
				userEntity.Update(userParams)
				topicRecordUpvote(topicEntity, userEntity, context.Get("ip").(string)) 
			}
		}
    	}

    	// Добавление в очередь
    	schedule.At(createFunc, scheduleContext, publicationTime, 0)

	return response.Send(w, 200, &response.Response{
		Data: nil,
		Meta: map[string]interface{}{
			"_notice": "Топик добавлен в очередь и скоро будет опубликован.",
		},
	})
}
