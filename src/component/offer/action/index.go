package offeraction

import (
	"fmt"
	"strings"
	"time"

	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/component/offer"
	"github.com/techfront/core/src/component/user"
	"github.com/techfront/core/src/lib/displayerror"
)

const (
	DEFAULT_LIMIT = 22
	DEFAULT_ORDER = "desc"
	DEFAULT_SORT = "id"
)

/**
* HandleIndex отображает список офферов.
 */
func HandleIndex(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Построение запроса
	query := offer.Published().Limit(DEFAULT_LIMIT)

	query.Where("score(offer_count_upvote, offer_count_downvote, offer_count_flag) > -6").Order("offer_created_at desc")

	// Фильтр по q (поисковой запрос)
	qFilter := params.Get("q")
	if len(qFilter) > 0 {
		qFilter = strings.Replace(qFilter, "_", "\\_", -1)
		qFilter = strings.Replace(qFilter, "%", "\\%", -1)
		wildcard := "%" + qFilter + "%"

		query.Where("offer_name ILIKE ? OR offer_url ILIKE ?", wildcard, wildcard)

		query.Order("rank(score(offer_count_upvote, offer_count_downvote, offer_count_flag), offer_created_at) desc, offer_id desc")
	}

	// Фильтр по u (id пользователя)
	uFilter := params.Get("u")
	if len(uFilter) > 0 {
		query.Where("offer_id_user=?", uFilter)
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

	return v.Render(context)
}

/**
* Функция задаёт базовые параметры для офферов.
*
* @param v *view.View Содержит инициализированый View.
* @param context router.Context Содержит интерфейс контекста.
 */
func setOfferMetadata(v *view.View, context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return err
	}

	// Определение следующей страницы
	nextPage := params.GetInt("page") + 1
	prevPage := params.GetInt("page") - 1

	// Задаем базовые параметры для офферов
	v.Vars["pub_date"] = time.Now()
	v.Vars["next_page"] = nextPage
	v.Vars["prev_page"] = prevPage

	v.Vars["next_page_link"] = fmt.Sprintf("?page=%d", nextPage)
	v.Vars["prev_page_link"] = fmt.Sprintf("?page=%d", prevPage)

	v.Vars["meta_title"] = "Офферы / Техфронт"
	v.Vars["meta_desc"] = "Сообщество энтузиастов, кому по нраву обсуждать технологии, исследования, стартапы и всё то, что происходит в сети."
	v.Vars["meta_keywords"] = "технологии, it, стартапы, проекты, новости, форум, сообщество, обсуждения"

	qFilter := params.Get("q")
	if len(qFilter) > 0 {
		v.Vars["next_page_link"] = fmt.Sprintf("?q=%s&page=%d", qFilter, nextPage)
		v.Vars["prev_page_link"] = fmt.Sprintf("?q=%s&page=%d", qFilter, prevPage)
		v.Vars["page_title"] = "Поиск по «" + qFilter + "»"
		v.Vars["meta_title"] = "Поиск по «" + qFilter + "» / Техфронт"
	}

	// Если определен параметр u, то меняем название на пользовательское
	uFilter := params.GetInt("u")
	if uFilter > 0 {
		userEntity, err := user.Find(uFilter)
		if err != nil {
			return router.InternalError(err, displayerror.PageNotFound...)
		}

		v.Vars["page_title"] = fmt.Sprintf("Все офферы %s", userEntity.Name)
		v.Vars["page_icon"] = "icon-paper-plane"
		v.Vars["meta_title"] = fmt.Sprintf("Все офферы %s / Техфронт", userEntity.Name)
	}

	path := strings.Replace(context.Path(), "/xml/", "", 1)
	if path == "/" {
		path = "/offers"
	}

	query := context.Request().URL.RawQuery
	if len(query) > 0 {
		query = "?" + query
	}

	url := fmt.Sprintf("/xml%s%s", path, query)
	v.Vars["meta_rss"] = url

	return nil
}
