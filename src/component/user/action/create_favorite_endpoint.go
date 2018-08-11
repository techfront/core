package useraction

import (
	"github.com/techfront/core/src/component/user"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/kernel/utils"
	"github.com/techfront/core/src/kernel/response"
)

/**
* HandleCreateFavoriteEndpoint добавляет топик в избранное.
 */
func HandleCreateFavoriteEndpoint(context router.Context) error {

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
		return response.Send(w, 404, response.NotFound("Что бы добавить топик в избранное необходимо авторизоваться."))
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

	// Построение параметров
	favoriteParams := map[string]string{
		"user_favorite_id_user": utils.ToString(userEntity.Id),
		"user_favorite_id_topic": utils.ToString(topicId),
	}

	// Проверка избранного
	if userEntity.IsTopicFavorited(topicId) {
		return response.Send(w, 401, response.InternalError("Топик уже был добавлен в избранное."))
	}

	_, err = user.CreateFavorite(favoriteParams)
	if err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}

	return response.Send(w, 200, &response.Response{
		Data: nil,
		Meta: map[string]interface{}{
			"_notice": "Топик успешно добавлен в избранное.",
		},
	})
}
