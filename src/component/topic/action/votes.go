package topicaction

import (
	"fmt"
	"log"

	"github.com/fragmenta/query"

	"github.com/techfront/core/src/kernel/router"

	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"

	"github.com/techfront/core/src/component/topic"
	"github.com/techfront/core/src/component/user"
)

/**
* HandleFlag обрабатывает такие маршруты как /topics/123/flag
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

	// Получение топика
	topicEntity, err := topic.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Получения пользователя
	userEntity := authorise.CurrentUser(context)
	ip := context.ClientIP()

	// Проверка прав
	if topicHasUserFlag(topicEntity, userEntity) {
		return router.NotAuthorizedError(err, "Flag Failed", "Sorry you are not allowed to flag twice, nice try!")
	}

	// Проверка прав на действие
	if !userEntity.CanFlag() {
		return router.NotAuthorizedError(err, "Flag Failed", "Sorry, you can't flag yet")
	}
	err = authorise.Resource(context, topicEntity)
	if err != nil {
		return router.NotAuthorizedError(err, "Flag Failed", "Sorry you are not allowed to flag")
	}

	// Обновляем колличество голосов
	err = topicEntity.Update(map[string]string{"topic_count_flag": fmt.Sprintf("%d", topicEntity.FlagCount+1.0)})
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Обновляем рейтинг пользователя создавшего топик
	topicUser, err := user.Find(topicEntity.UserId)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}
	err = adjustUserScore(topicUser, userEntity.ComputeDeltaScore()*(-5.0))
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Записываем данные в отдельную таблицу
	if _, err := query.Exec("INSERT INTO tf_flag (flag_created_at, flag_id_topic, flag_id_user, flag_ip) VALUES(now(),$1,$2,$3)", topicEntity.Id, userEntity.Id, ip); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return nil
}

/**
* HandleDownvote handles POST to /topics/123/downvote
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

	// Получение топика
	topicEntity, err := topic.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Получение пользоватея
	userEntity := authorise.CurrentUser(context)
	ip := context.ClientIP()

	// Проверка прав
	if !userEntity.Admin() {
		if topicHasUserVote(topicEntity, userEntity) {
			return router.NotAuthorizedError(err, "Vote Failed", "Sorry you are not allowed to vote twice, nice try!")
		}
	}

	// Проверка прав не действие
	if !userEntity.CanDownvote() {
		return router.NotAuthorizedError(err, "Vote Failed", "Sorry, you can't downvote yet")
	}
	err = authorise.Resource(context, topicEntity)
	if err != nil {
		return router.NotAuthorizedError(err, "Vote Failed", "Sorry you are not allowed to vote")
	}

	// Обновляем колличество голосов
	err = topicEntity.Update(map[string]string{"topic_count_downvote": fmt.Sprintf("%d", topicEntity.DownvoteCount+1.0)})
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Обновляем рейтинг пользователя создавшего топик
	topicUser, err := user.Find(topicEntity.UserId)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}
	err = adjustUserScore(topicUser, userEntity.ComputeDeltaScore()*(-1.0))
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Записываем данные в отдельную таблицу
	if _, err := query.Exec("INSERT INTO tf_downvote (downvote_created_at, downvote_id_topic, downvote_id_user, downvote_ip) VALUES(now(),$1,$2,$3)", topicEntity.Id, userEntity.Id, ip); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return nil
}

/**
* HandleUpvote handles POST to /topics/123/upvote
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

	// Получение топика
	topicEntity, err := topic.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Получение пользователя
	userEntity := authorise.CurrentUser(context)
	ip := context.ClientIP()

	// Проверка прав на действие
	if !userEntity.Admin() {
		if topicHasUserVote(topicEntity, userEntity) {
			return router.NotAuthorizedError(err, "Vote Failed", "Sorry you are not allowed to vote twice, nice try!")
		}
	}

	// Проверка прав на действие
	if !userEntity.CanUpvote() {
		return router.NotAuthorizedError(err, "Vote Failed", "Sorry, you can't upvote yet")
	}

	err = authorise.Resource(context, topicEntity)
	if err != nil {
		return router.NotAuthorizedError(err, "Vote Failed", "Sorry you are not allowed to vote")
	}

	// Обновляем колличество голосов
	err = topicEntity.Update(map[string]string{"topic_count_upvote": fmt.Sprintf("%d", topicEntity.UpvoteCount+1.0)})
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Обновляем рейтинг и силу пользователя создавшего топик
	topicUser, err := user.Find(topicEntity.UserId)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}
	err = adjustUserPower(topicUser, userEntity.ComputeDeltaPower()*(+1.0))
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}
	err = adjustUserScore(topicUser, userEntity.ComputeDeltaScore()*(+1.0))
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Записываем данные в отдельную таблицу
	if _, err := query.Exec("INSERT INTO tf_upvote (upvote_created_at, upvote_id_topic, upvote_id_user, upvote_ip) VALUES(now(),$1,$2,$3)", topicEntity.Id, userEntity.Id, ip); err != nil {
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
* Функция upvoteTopic добавляет один положительный голос без проверок
 */
func topicUpvote(topicEntity *topic.Topic, userEntity *user.User, ip string) error {

	// Обновляем колличество голосов
	err := topicEntity.Update(map[string]string{"topic_count_upvote": fmt.Sprintf("%d", topicEntity.UpvoteCount+1.0)})
	if err != nil {
		return err
	}

	// Обновляем рейтинг и силу пользователя создавшего топик
	topicUser, err := user.Find(topicEntity.UserId)
	if err != nil {
		return err
	}
	err = adjustUserPower(topicUser, userEntity.ComputeDeltaPower()*(+1.0))
	if err != nil {
		return err
	}
	err = adjustUserScore(topicUser, userEntity.ComputeDeltaScore()*(+1.0))
	if err != nil {
		return err
	}

	return topicRecordUpvote(topicEntity, userEntity, ip)
}

/**
* Функция recordUpvote записывает данные о голосе в отдельную таблицу.
 */
func topicRecordUpvote(topicEntity *topic.Topic, userEntity *user.User, ip string) error {
	if _, err := query.Exec("INSERT INTO tf_upvote (upvote_created_at, upvote_id_topic, upvote_id_user, upvote_ip) VALUES(now(),$1,$2,$3)", topicEntity.Id, userEntity.Id, ip); err != nil {
		return err
	}

	return nil
}

/**
* topicHasUserVote возвращает true если пользователь голосовал за топик.
 */
func topicHasUserVote(topicEntity *topic.Topic, userEntity *user.User) bool {

	// Проверка Upvote
	if results, err := query.New("tf_upvote", "upvote_id_topic").Where("upvote_id_topic=?", topicEntity.Id).Where("upvote_id_user=?", userEntity.Id).Results(); err == nil && len(results) == 0 {
		return false
	}

	// Проверка Downvote
	if results, err := query.New("tf_downvote", "downvote_id_topic").Where("downvote_id_topic=?", topicEntity.Id).Where("downvote_id_user=?", userEntity.Id).Results(); err == nil && len(results) == 0 {
		return false
	}

	return true
}

/**
* topicHasUserFlag возвращает true если пользователь уже отправил жалобу.
 */
func topicHasUserFlag(topicEntity *topic.Topic, userEntity *user.User) bool {

	// Проверка Flag
	if results, err := query.New("tf_flag", "flag_id_topic").Where("flag_id_topic=?", topicEntity.Id).Where("flag_id_user=?", userEntity.Id).Results(); err == nil && len(results) == 0 {
		return false
	}

	return true
}
