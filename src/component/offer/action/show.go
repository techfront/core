package offeraction

import (
	"fmt"
	"strings"

	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/cache"
	"github.com/techfront/core/src/lib/displayerror"

	"github.com/techfront/core/src/component/comment"
	"github.com/techfront/core/src/component/offer"
)

/**
* HandleShow displays a single offer.
 */
func HandleShow(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение оффера
	offerEntity, err := offer.Find(params.GetInt("id"))
	if err != nil {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Скрываем оффер если он отклонен и пользователь не админ.
	if offerEntity.IsRejected() && !authorise.CurrentUser(context).Admin() {
		return router.NotFoundError(err, displayerror.PageNotFound...)
	}

	// Редирект на канонический URL
	if context.Path() != offerEntity.URLShow() {
		return router.Redirect(context, offerEntity.URLShow())
	}

	// Получение интересных офферов для виджета
	// Оптимизация: кэшированние выборки
	relatedOfferKey := fmt.Sprintf("offer__id-%d--related", params.GetInt("id"))

	var relatedOfferList []*offer.Offer

	if err := cache.Get(relatedOfferKey, &relatedOfferList); err != nil {
		relatedOfferQuery := offerEntity.Related().Limit(8)

		results, err := offer.FindAll(relatedOfferQuery)
		if err != nil {
			return router.InternalError(err, displayerror.UnknownError...)
		}

		relatedOfferList = results

		if err := cache.Set(relatedOfferKey, relatedOfferList, 360); err != nil {
			return router.InternalError(err, displayerror.UnknownError...)
		}
	}

	// Получение комментариев
	commentsQuery := comment.Where("comment_id_offer=?", offerEntity.Id).Order(comment.RANK_ORDER)
	commentList, err := comment.FindAllWithChild(commentsQuery)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Отображаение шаблона
	v := view.New(
		"component/offer/template/show",
		"component/offer/template/row_mini",
		"component/offer/template/widget/widget_related",
		"component/offer/template/widget/widget_share",
		"component/comment/template/comment",
		"component/comment/template/form_embed",
	)

	switch params.Get("message") {
	case "info--offer_in_the_queue":
		v.Vars["message_type"] = "default"
		v.Vars["message"] = "Спасибо, оффер успешно создан и ожидает свою очередь."
	case "error--comment_invalid_create_timeout":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Добавить новый комментарий можно через 1 минуту. Подождите и попробуйте снова."
	case "error--comment_invalid_length":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Комментарий должен быть длинее 10 символов."
	}

	var desc string
	if len(offerEntity.Text) > 155 {
		desc = offerEntity.Text[:155]
		lastWordLenght := len(desc) - strings.LastIndex(desc, " ")
		desc = fmt.Sprintf("%s...", desc[:155 - lastWordLenght])
	}

	v.Vars["offer"] = offerEntity
	v.Vars["meta_title"] = offerEntity.Name + " / Техфронт"
	v.Vars["meta_desc"] = desc
	v.Vars["meta_keywords"] = ""
	v.Vars["comments"] = commentList
	v.Vars["related_offers"] = relatedOfferList

	return v.Render(context)
}
