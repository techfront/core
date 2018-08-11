package offeraction

import (
	"strings"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"
	"github.com/techfront/core/src/component/offer"
	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleFavorites отображает список офферов из закладок.
 */
func HandleFavorited(context router.Context) error {
	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение пользователя
	userEntity := authorise.CurrentUser(context)
	if userEntity.Id == 0 {
		return router.NotAuthorizedError(nil, displayerror.AccessDenied...)
	}

	// Построение запроса
	query := offer.Favorited(userEntity.Id).Limit(DEFAULT_LIMIT)

	// Фильтр по q (поисковой запрос)
	qFilter := params.Get("q")
	if len(qFilter) > 0 {
		qFilter = strings.Replace(qFilter, "_", "\\_", -1)
		qFilter = strings.Replace(qFilter, "%", "\\%", -1)
		wildcard := "%" + qFilter + "%"

		query.Where("offer_name ILIKE ? OR offer_url ILIKE ?", wildcard, wildcard)
	}

	// Определяем текущую страницу и делаем отступ
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
		"component/offer/template/index",
		"component/offer/template/row_mini",
	)

	v.Vars["page"] = currentPage
	v.Vars["offers"] = results
	v.Vars["offers_format"] = -1
	v.Vars["offers_count"] = len(results)

	if err := setOfferMetadata(v, context); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	v.Vars["page_icon"] = "icon-bookmark"
	v.Vars["page_title"] = "Избранные офферы"
	v.Vars["meta_title"] = "Избранные офферы / Техфронт"

	return v.Render(context)
}