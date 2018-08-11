package offeraction

import (
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/component/offer"
	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleHome displays a list of offers using gravity to order them
* used for the home page for gravity rank see votes.go
* responds to GET /
 */
func HandleHome(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Построение запроса
	query := offer.Published().Limit(DEFAULT_LIMIT)

	query.Where("offer_count_upvote > 0")

	query.Order("offer_status desc, rank(score(offer_count_upvote, offer_count_downvote, offer_count_flag), offer_created_at) desc, offer_id desc")

	currentPage := int(params.GetInt("page"))
	if currentPage > 0 {
		query.Offset(DEFAULT_LIMIT * currentPage)
	}

	// Получение запроса
	results, err := offer.FindAll(query)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Отображение шаблона
	v := view.New(
		"component/offer/template/home",
		"component/offer/template/row_mini",
	)

	v.Vars["offers"] = results

	if err := setOfferMetadata(v, context); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return v.Render(context)
}
