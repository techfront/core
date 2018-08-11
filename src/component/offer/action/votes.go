package offeraction

import (
	"fmt"
	"log"

	"github.com/fragmenta/query"

	"github.com/techfront/core/src/kernel/router"

	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"

	"github.com/techfront/core/src/component/offer"
	"github.com/techfront/core/src/component/user"
)

/**
* HandleFlag обрабатывает такие маршруты как /offers/123/flag
* Только для Администратора.
 */
func HandleFlag(context router.Context) error {

	// Проверка прав
	err := authorise.AuthenticityToken(context)
	if err != nil {
		return router.NotAuthorizedError(err, "Flag Failed", "CSRF failure")
	}

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение оффера
	offerEntity, err := offer.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Получения пользователя
	userEntity := authorise.CurrentUser(context)
	ip := context.ClientIP()

	// Проверка прав
	if offerHasUserFlag(offerEntity, userEntity) {
		return router.NotAuthorizedError(err, "Flag Failed", "Sorry you are not allowed to flag twice, nice try!")
	}

	// Проверка прав на действие
	if !userEntity.CanFlag() {
		return router.NotAuthorizedError(err, "Flag Failed", "Sorry, you can't flag yet")
	}
	err = authorise.Resource(context, offerEntity)
	if err != nil {
		return router.NotAuthorizedError(err, "Flag Failed", "Sorry you are not allowed to flag")
	}

	// Обновляем колличество голосов
	err = offerEntity.Update(map[string]string{"offer_count_flag": fmt.Sprintf("%d", offerEntity.FlagCount+1.0)})
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Обновляем рейтинг пользователя создавшего оффер
	offerUser, err := user.Find(offerEntity.UserId)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}
	err = adjustUserScore(offerUser, userEntity.ComputeDeltaScore()*(-5.0))
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Записываем данные в отдельную таблицу
	if _, err := query.Exec("INSERT INTO tf_flag (flag_created_at, flag_id_offer, flag_id_user, flag_ip) VALUES(now(),$1,$2,$3)", offerEntity.Id, userEntity.Id, ip); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return nil
}

/**
* HandleDownvote handles POST to /offers/123/downvote
 */
func HandleDownvote(context router.Context) error {

	// Проверка прав
	err := authorise.AuthenticityToken(context)
	if err != nil {
		return router.NotAuthorizedError(err, "Downvote Failed", "CSRF failure")
	}

	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение оффера
	offerEntity, err := offer.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Получение пользоватея
	userEntity := authorise.CurrentUser(context)
	ip := context.ClientIP()

	// Проверка прав
	if !userEntity.Admin() {
		if offerHasUserVote(offerEntity, userEntity) {
			return router.NotAuthorizedError(err, "Vote Failed", "Sorry you are not allowed to vote twice, nice try!")
		}
	}

	// Проверка прав не действие
	if !userEntity.CanDownvote() {
		return router.NotAuthorizedError(err, "Vote Failed", "Sorry, you can't downvote yet")
	}
	err = authorise.Resource(context, offerEntity)
	if err != nil {
		return router.NotAuthorizedError(err, "Vote Failed", "Sorry you are not allowed to vote")
	}

	// Обновляем колличество голосов
	err = offerEntity.Update(map[string]string{"offer_count_downvote": fmt.Sprintf("%d", offerEntity.DownvoteCount+1.0)})
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Обновляем рейтинг пользователя создавшего оффер
	offerUser, err := user.Find(offerEntity.UserId)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}
	err = adjustUserScore(offerUser, userEntity.ComputeDeltaScore()*(-1.0))
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Записываем данные в отдельную таблицу
	if _, err := query.Exec("INSERT INTO tf_downvote (downvote_created_at, downvote_id_offer, downvote_id_user, downvote_ip) VALUES(now(),$1,$2,$3)", offerEntity.Id, userEntity.Id, ip); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return nil
}

/**
* HandleUpvote handles POST to /offers/123/upvote
 */
func HandleUpvote(context router.Context) error {

	// Проверка прав
	err := authorise.AuthenticityToken(context)
	if err != nil {
		return router.NotAuthorizedError(err, "Upvote Failed", "CSRF failure")
	}

	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение оффера
	offerEntity, err := offer.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Получение пользователя
	userEntity := authorise.CurrentUser(context)
	ip := context.ClientIP()

	// Проверка прав на действие
	if !userEntity.Admin() {
		if offerHasUserVote(offerEntity, userEntity) {
			return router.NotAuthorizedError(err, "Vote Failed", "Sorry you are not allowed to vote twice, nice try!")
		}
	}

	// Проверка прав на действие
	if !userEntity.CanUpvote() {
		return router.NotAuthorizedError(err, "Vote Failed", "Sorry, you can't upvote yet")
	}

	err = authorise.Resource(context, offerEntity)
	if err != nil {
		return router.NotAuthorizedError(err, "Vote Failed", "Sorry you are not allowed to vote")
	}

	// Обновляем колличество голосов
	err = offerEntity.Update(map[string]string{"offer_count_upvote": fmt.Sprintf("%d", offerEntity.UpvoteCount+1.0)})
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Обновляем рейтинг и силу пользователя создавшего оффер
	offerUser, err := user.Find(offerEntity.UserId)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}
	err = adjustUserPower(offerUser, userEntity.ComputeDeltaPower()*(+1.0))
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}
	err = adjustUserScore(offerUser, userEntity.ComputeDeltaScore()*(+1.0))
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Записываем данные в отдельную таблицу
	if _, err := query.Exec("INSERT INTO tf_upvote (upvote_created_at, upvote_id_offer, upvote_id_user, upvote_ip) VALUES(now(),$1,$2,$3)", offerEntity.Id, userEntity.Id, ip); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return nil
}

/**
* removeUserScore обновляет рейтинг пользователя
 */
func adjustUserScore(userEntity *user.User, delta float64) error {

	// Обновляем рейтинг пользователя
	err := userEntity.Update(map[string]string{"user_score": fmt.Sprintf("%f", userEntity.Score+delta)})
	if err != nil {
		log.Print(err)
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return nil
}

/**
* removeUserPower обновляет силу пользователя
 */
func adjustUserPower(userEntity *user.User, delta float64) error {

	// Обновляем силу пользователя
	err := userEntity.Update(map[string]string{"user_power": fmt.Sprintf("%f", userEntity.Power+delta)})
	if err != nil {
		log.Print(err)
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return nil
}

/**
* Функция upvoteOffer добавляет один положительный голос без проверок
 */
func offerUpvote(offerEntity *offer.Offer, userEntity *user.User, ip string) error {

	// Обновляем колличество голосов
	err := offerEntity.Update(map[string]string{"offer_count_upvote": fmt.Sprintf("%d", offerEntity.UpvoteCount+1.0)})
	if err != nil {
		return err
	}

	// Обновляем рейтинг и силу пользователя создавшего оффер
	offerUser, err := user.Find(offerEntity.UserId)
	if err != nil {
		return err
	}
	err = adjustUserPower(offerUser, userEntity.ComputeDeltaPower()*(+1.0))
	if err != nil {
		return err
	}
	err = adjustUserScore(offerUser, userEntity.ComputeDeltaScore()*(+1.0))
	if err != nil {
		return err
	}

	return offerRecordUpvote(offerEntity, userEntity, ip)
}

/**
* Функция recordUpvote записывает данные о голосе в отдельную таблицу.
 */
func offerRecordUpvote(offerEntity *offer.Offer, userEntity *user.User, ip string) error {
	if _, err := query.Exec("INSERT INTO tf_upvote (upvote_created_at, upvote_id_offer, upvote_id_user, upvote_ip) VALUES(now(),$1,$2,$3)", offerEntity.Id, userEntity.Id, ip); err != nil {
		return err
	}

	return nil
}

/**
* offerHasUserVote возвращает true если пользователь голосовал за оффер.
 */
func offerHasUserVote(offerEntity *offer.Offer, userEntity *user.User) bool {

	// Проверка Upvote
	if results, err := query.New("tf_upvote", "upvote_id_offer").Where("upvote_id_offer=?", offerEntity.Id).Where("upvote_id_user=?", userEntity.Id).Results(); err == nil && len(results) == 0 {
		return false
	}

	// Проверка Downvote
	if results, err := query.New("tf_downvote", "downvote_id_offer").Where("downvote_id_offer=?", offerEntity.Id).Where("downvote_id_user=?", userEntity.Id).Results(); err == nil && len(results) == 0 {
		return false
	}

	return true
}

/**
* offerHasUserFlag возвращает true если пользователь уже отправил жалобу.
 */
func offerHasUserFlag(offerEntity *offer.Offer, userEntity *user.User) bool {

	// Проверка Flag
	if results, err := query.New("tf_flag", "flag_id_offer").Where("flag_id_offer=?", offerEntity.Id).Where("flag_id_user=?", userEntity.Id).Results(); err == nil && len(results) == 0 {
		return false
	}

	return true
}
