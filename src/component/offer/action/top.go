package offeraction

import (
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	// "github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"

	"github.com/techfront/core/src/component/offer"
)

/**
* HandleTop displays a list of offers the user has top in the past
 */
func HandleTop(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Построение запроса
	query := offer.Published().Limit(DEFAULT_LIMIT)

	query.Where("score(offer_count_upvote, offer_count_downvote, offer_count_flag) > 0")
	query.Order("rank(score(offer_count_upvote, offer_count_downvote, offer_count_flag), offer_created_at) desc, score(offer_count_upvote, offer_count_downvote, offer_count_flag) desc, offer_id desc")

	/*
		// Выбрать для авторизованного пользователя
		userEntity := authorise.CurrentUser(context)
		if !user.Anon() {
			// Can we use a join instead?
			v := query.New("votes", "offer_id").Select("select offer_id as id from votes").Where("user_id=? AND offer_id IS NOT NULL AND points > 0", u.Id)

			offerIDs := v.ResultIDs()
			if len(offerIDs) > 0 {
				query.WhereIn("id", offerIDs)
			}
		}
	*/

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
	v.Vars["offers_format"] = -1
	v.Vars["offers_count"] = len(results)

	if params.Get("format") == ".xml" {
		v.Extension = "xml.got"
		v.Templates = []string{"component/offer/template/index"}
	}

	if err := setOfferMetadata(v, context); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	v.Vars["meta_title"] =  "Топ офферов / Техфронт"

	return v.Render(context)
}
