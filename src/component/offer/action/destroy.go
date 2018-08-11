package offeraction

import (
	"fmt"
	"github.com/techfront/core/src/kernel/router"

	"github.com/techfront/core/src/component/offer"
	"github.com/techfront/core/src/component/user"

	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleDestroy handles a DESTROY request for offers
 */
func HandleDestroy(context router.Context) error {
	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение оффера
	offerEntity, err := offer.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Получения пользователя создавшего оффер
	userEntity, err := user.Find(offerEntity.UserId)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Если оффер отклонен, то уже удалить нельзя. Если пользователь не админ
	if offerEntity.IsRejected() && !authorise.CurrentUser(context).Admin() {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Проверка прав
	err = authorise.Resource(context, offerEntity)
	if err != nil {
		return router.NotAuthorizedError(err, displayerror.AccessDenied...)
	}

	// Удаление оффера
	err = offerEntity.Destroy()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Обновление колличества офферов пользователя
	userParams := map[string]string{"user_count_offer": fmt.Sprintf("%d", userEntity.OfferCount-1)}
	err = userEntity.Update(userParams)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return router.Redirect(context, "/")
}
