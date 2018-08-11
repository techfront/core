package topicaction

import (
	"fmt"
	"github.com/techfront/core/src/kernel/router"

	"github.com/techfront/core/src/component/topic"
	"github.com/techfront/core/src/component/user"

	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleDestroy handles a DESTROY request for topics
 */
func HandleDestroy(context router.Context) error {
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

	// Получения пользователя создавшего топик
	userEntity, err := user.Find(topicEntity.UserId)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Если топик отклонен, то уже удалить нельзя. Если пользователь не админ
	if topicEntity.IsRejected() && !authorise.CurrentUser(context).Admin() {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Проверка прав
	err = authorise.Resource(context, topicEntity)
	if err != nil {
		return router.NotAuthorizedError(err, displayerror.AccessDenied...)
	}

	// Удаление топика
	err = topicEntity.Destroy()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Обновление колличества топиков пользователя
	userParams := map[string]string{"user_count_topic": fmt.Sprintf("%d", userEntity.TopicCount-1)}
	err = userEntity.Update(userParams)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return router.Redirect(context, "/")
}
