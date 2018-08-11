package useraction

import (
	"github.com/techfront/core/src/component/user"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/kernel/response"
)

/**
* HandleDestroyFavoriteEndpoint удаляет топик из избранного.
 */
func HandleDestroyFavoriteEndpoint(context router.Context) error {

	w := context.Writer()

	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}

	topicId := params.GetInt("id_topic")
	userId := params.GetInt("id_user")

	// Проверка идентификатора пользователя
	if userId == 0 {
		return response.Send(w, 404, response.NotFound("Что бы удалить топик из избранного необходимо авторизоваться."))
	}

	// Получение пользователя
	userEntity, err := user.Find(userId)
	if err != nil {
		return response.Send(w, 404, response.NotFound("Ошибка 404. Ресурс не найден."))
	}

	// Проверка прав
	err = authorise.ResourceAndAuthenticity(context, userEntity)
	if err != nil {
		return response.Send(w, 401, response.Unauthorized("Ошибка 401. Необходима авторизация."))
	}

	// Получение закладки
	favoriteEntity, err := userEntity.FindFavoriteByTopicId(topicId)
	if err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}

	// Удаление закладки
	err = favoriteEntity.DestroyFavorite()
	if err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}

	return response.Send(w, 200, &response.Response{
		Data: nil,
		Meta: map[string]interface{}{
			"_notice": "Топик успешно удален из избранного.",
		},
	})
}
