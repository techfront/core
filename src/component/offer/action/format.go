package offeraction

import (
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/component/offer"
	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleFormat отображает список офферов на странице формата
 */
func HandleFormat(context router.Context) error {
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Построение запроса
	query := offer.Published().Where("score(offer_count_upvote, offer_count_downvote, offer_count_flag) > -6")

	query.Order("rank(score(offer_count_upvote, offer_count_downvote, offer_count_flag), offer_created_at) desc, score(offer_count_upvote, offer_count_downvote, offer_count_flag) desc, offer_id desc").Limit(DEFAULT_LIMIT)

	var formatTitle string
	var formatID int64

	// Все форматы
	formatID = -1

	// Определение конкретного формата с помощью параметра в URL
	switch params.Get("format") {
	case "offer":
		formatTitle = "офферы"
		formatID = offer.FormatOffer
		query.Where("offer_id_format = ? OR offer_id_format IS NULL", offer.FormatOffer)
	case "job":
		formatTitle = "вакансии"
		formatID = offer.FormatJob
		query.Where("offer_id_format = ?", offer.FormatJob)
	default:
		return router.NotFoundError(nil, displayerror.PageNotFound...)
	}

	// Определение текущей страницы и добавление отступа
	currentPage := int(params.GetInt("page"))
	if currentPage > 0 {
		query.Offset(DEFAULT_LIMIT * currentPage)
	}

	// Получение офферов
	results, err := offer.FindAll(query)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Отображение шаблона
	v := view.New(
		"component/offer/template/index",
		"component/offer/template/row_mini",
	)

	v.Vars["page"] = currentPage
	v.Vars["offers"] = results
	v.Vars["offers_count"] = len(results)
	v.Vars["offers_format"] = formatID

	if err := setOfferMetadata(v, context); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	v.Vars["page_title"] = "Все " + formatTitle
	v.Vars["meta_title"] = "Все " + formatTitle + " / Техфронт"

	return v.Render(context)
}
