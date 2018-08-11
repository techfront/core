package commentaction

import (
	"fmt"
	"github.com/fragmenta/query"

	"github.com/techfront/core/src/kernel/router"

	"github.com/techfront/core/src/component/comment"
	"github.com/techfront/core/src/component/user"

	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleFlag обрабатывает такие маршруты как /comments/123/flag
* Только для Администратора.
 */
func HandleFlag(context router.Context) error {

	// Проверка прав
	err := authorise.AuthenticityToken(context)
	if err != nil {
		return router.NotAuthorizedError(err, "Flag Failed", "CSRF failure")
	}

	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение комментария
	commentEntity, err := comment.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Получение пользователя
	userEntity := authorise.CurrentUser(context)
	ip := context.ClientIP()

	// Проверка прав на действие
	if commentHasUserFlag(commentEntity, userEntity) {
		return router.NotAuthorizedError(err, "Flag Failed", "Sorry you are not allowed to flag twice, nice try!")
	}

	// Проверка прав на действие
	if !userEntity.CanFlag() {
		return router.NotAuthorizedError(err, "Flag Failed", "Sorry, you can't flag yet")
	}

	// Обновляем колличество голосов
	err = commentEntity.Update(map[string]string{"comment_count_flag": fmt.Sprintf("%d", commentEntity.FlagCount+1.0)})
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Обновляем рейтинг пользователя создавшего топик
	commentUser, err := user.Find(commentEntity.UserId)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}
	err = adjustUserScore(commentUser, userEntity.ComputeDeltaScore()*(-5.0))
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Записываем данные в отдельную таблицу
	if _, err := query.Exec("INSERT INTO tf_flag (flag_created_at, flag_id_comment, flag_id_user, flag_ip) VALUES(now(),$1,$2,$3)", commentEntity.Id, userEntity.Id, ip); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return nil
}

/**
* HandleDownvote handles POST to /comments/123/downvote
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

	// Получение комментария
	commentEntity, err := comment.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Получение пользователя
	userEntity := authorise.CurrentUser(context)
	ip := context.ClientIP()

	// Проверка прав на действие
	if !userEntity.Admin() {
		if commentHasUserVote(commentEntity, userEntity) {
			return router.NotAuthorizedError(err, "Vote Failed", "Sorry you are not allowed to vote twice, nice try!")
		}
	}

	// Проверка прав на действие
	if !userEntity.CanDownvote() {
		return router.NotAuthorizedError(err, "Vote Failed", "Sorry, you can't downvote yet")
	}
	err = authorise.Resource(context, commentEntity)
	if err != nil {
		return router.NotAuthorizedError(err, "Vote Failed", "Sorry you are not allowed to vote")
	}

	// Обновляем колличество голосов
	err = commentEntity.Update(map[string]string{"comment_count_downvote": fmt.Sprintf("%d", commentEntity.DownvoteCount+1.0)})
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Обновляем рейтинг пользователя создавшего топик
	commentUser, err := user.Find(commentEntity.UserId)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}
	err = adjustUserScore(commentUser, userEntity.ComputeDeltaScore()*(-1.0))
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Записываем данные в отдельную таблицу
	if _, err := query.Exec("INSERT INTO tf_downvote (downvote_created_at, downvote_id_comment, downvote_id_user, downvote_ip) VALUES(now(),$1,$2,$3)", commentEntity.Id, userEntity.Id, ip); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return nil
}

/**
* HandleUpvote handles POST to /comments/123/upvote
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

	// Получение комментария
	commentEntity, err := comment.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Получение пользователя
	userEntity := authorise.CurrentUser(context)
	ip := context.ClientIP()

	if !userEntity.Admin() {
		if commentHasUserVote(commentEntity, userEntity) {
			return router.NotAuthorizedError(err, "Vote Failed", "Sorry you are not allowed to vote twice, nice try!")
		}
	}

	// Проверка прав на действие
	if !userEntity.CanUpvote() {
		return router.NotAuthorizedError(err, "Vote Failed", "Sorry, you can't upvote yet")
	}
	err = authorise.Resource(context, commentEntity)
	if err != nil {
		return router.NotAuthorizedError(err, "Vote Failed", "Sorry you are not allowed to vote")
	}

	// Обновляем колличество голосов
	err = commentEntity.Update(map[string]string{"comment_count_upvote": fmt.Sprintf("%d", commentEntity.UpvoteCount+1.0)})
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Обновляем рейтинг пользователя создавшего топик
	commentUser, err := user.Find(commentEntity.UserId)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}
	err = adjustUserPower(commentUser, userEntity.ComputeDeltaPower()*(+1.0))
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}
	err = adjustUserScore(commentUser, userEntity.ComputeDeltaScore()*(+1.0))
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Записываем данные в отдельную таблицу
	if _, err := query.Exec("INSERT INTO tf_upvote (upvote_created_at, upvote_id_comment, upvote_id_user, upvote_ip) VALUES(now(),$1,$2,$3)", commentEntity.Id, userEntity.Id, ip); err != nil {
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
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return nil
}

/**
* Функция commentUpvote добавляет один положительный голос без проверок
 */
func commentUpvote(commentEntity *comment.Comment, userEntity *user.User, ip string) error {

	// Обновляем колличество голосов
	err := commentEntity.Update(map[string]string{"comment_count_upvote": fmt.Sprintf("%d", commentEntity.UpvoteCount+1.0)})
	if err != nil {
		return err
	}

	// Обновляем рейтинг и силу пользователя создавшего коммент
	commentUser, err := user.Find(commentEntity.UserId)
	if err != nil {
		return err
	}
	err = adjustUserPower(commentUser, userEntity.ComputeDeltaPower()*(+1.0))
	if err != nil {
		return err
	}
	err = adjustUserScore(commentUser, userEntity.ComputeDeltaScore()*(+1.0))
	if err != nil {
		return err
	}

	return commentRecordUpvote(commentEntity, userEntity, ip)
}

/**
* Функция commentRecordUpvote записывает данные о голосе в отдельную таблицу.
 */
func commentRecordUpvote(commentEntity *comment.Comment, userEntity *user.User, ip string) error {
	if _, err := query.Exec("INSERT INTO tf_upvote (upvote_created_at, upvote_id_comment, upvote_id_user, upvote_ip) VALUES(now(),$1,$2,$3)", commentEntity.Id, userEntity.Id, ip); err != nil {
		return err
	}

	return nil
}

/**
* commentHasUserVote возвращает true если пользователь голосовал за коммент.
 */
func commentHasUserVote(commentEntity *comment.Comment, userEntity *user.User) bool {

	// Проверка Upvote
	if results, err := query.New("tf_upvote", "upvote_id_comment").Where("upvote_id_comment=?", commentEntity.Id).Where("upvote_id_user=?", userEntity.Id).Results(); err == nil && len(results) == 0 {
		return false
	}

	// Проверка Downvote
	if results, err := query.New("tf_downvote", "downvote_id_comment").Where("downvote_id_comment=?", commentEntity.Id).Where("downvote_id_user=?", userEntity.Id).Results(); err == nil && len(results) == 0 {
		return false
	}

	return true
}

/**
* commentHasUserFlag возвращает true если пользователь уже отправил жалобу.
 */
func commentHasUserFlag(commentEntity *comment.Comment, userEntity *user.User) bool {

	// Проверка Flag
	if results, err := query.New("tf_flags", "flag_id_comment").Where("flag_id_comment=?", commentEntity.Id).Where("flag_id_user=?", userEntity.Id).Results(); err == nil && len(results) == 0 {
		return false
	}

	return true
}
