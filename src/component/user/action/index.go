package useraction

import (
	"fmt"
	"strings"
	"github.com/techfront/core/src/component/user"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"
	"github.com/techfront/core/src/lib/displayerror"
)

const (
	DEFAULT_LIMIT = 100
	DEFAULT_ORDER = "desc"
)

/**
* HandleIndex serves a GET request at /users
 */
func HandleIndex(context router.Context) error {
	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Построение запроса
	query := user.Query()

	// Обработка q
	qFilter := params.Get("q")
	if len(qFilter) > 0 {
		qFilter = strings.Replace(qFilter, "_", "\\_", -1)
		qFilter = strings.Replace(qFilter, "%", "\\%", -1)
		wildcard := "%" + qFilter + "%"

		query.Where("user_name ILIKE ? OR user_fullname ILIKE ?", wildcard, wildcard)
	}

	// Обработка order
	order := params.Get("order")
	if order != "desc" && order != "asc" {
		order = DEFAULT_ORDER
	}
	
	// Обработка sort
	var currentSort string
	sort := params.Get("sort")
	switch sort {
		case "score":
			currentSort = "score"
			sort = fmt.Sprintf("user_score %s, user_name asc", order)
			query.Order(sort)
		case "count_topic":
			currentSort = "count_topic"
			sort = fmt.Sprintf("user_count_topic %s, user_name asc", order)
			query.Order(sort)
		case "count_comment":
			currentSort = "count_comment"
			sort = fmt.Sprintf("user_count_comment %s, user_name asc", order)
			query.Order(sort)
		case "power":
			currentSort = "power"
			sort = fmt.Sprintf("user_power %s, user_name asc", order)
			query.Order(sort)
		case "name":
			currentSort = "name"
			sort = fmt.Sprintf("user_name %s, user_id desc", order)
			query.Order(sort)
		default:
			currentSort = "score"
			sort = fmt.Sprintf("user_score %s, user_name asc", order)
			query.Order(sort)
	}

	// Обработка page
	currentPage := int(params.GetInt("page"))
	if currentPage > 0 {
		query.Offset(DEFAULT_LIMIT * currentPage)
	}

	query.Limit(DEFAULT_LIMIT)

	// Получение запроса
	userList, err := user.FindAll(query)
	if err != nil {
		return router.InternalError(nil, displayerror.UnknownError...)
	}

	// Получения колличества пользователей
	userCount, err := user.Query().Count()
	if err != nil {
		return router.InternalError(nil, displayerror.UnknownError...)
	}

	// Получение колличества админов
	adminCount, err := user.Query().Where("user_role=100").Count()
	if err != nil {
		return router.InternalError(nil, displayerror.UnknownError...)
	}

	// Отображение шаблона
	v := view.New("component/user/template/index", "component/user/template/row")

	nextPage := params.GetInt("page") + 1
	prevPage := params.GetInt("page") - 1

	v.Vars["next_page"] = nextPage
	v.Vars["prev_page"] = prevPage

	r := context.Request()
	currentQuery := r.URL.Query()

	currentQuery.Set("page", fmt.Sprintf("%d", nextPage))
	v.Vars["next_page_link"] = "?" + currentQuery.Encode()

	currentQuery.Set("page", fmt.Sprintf("%d", prevPage))
	v.Vars["prev_page_link"] = "?" + currentQuery.Encode()

	v.Vars["current_sort"] = currentSort
	v.Vars["current_order"] = order

	v.Vars["users"] = userList
	v.Vars["page_title"] = "Сообщество"
	v.Vars["page_icon"] = "icon-globe-1"
	v.Vars["meta_title"] = "Сообщество / Техфронт"
	v.Vars["meta_desc"] = "Сообщество энтузиастов, кому по нраву обсуждать технологии, исследования, стартапы и всё то, что происходит в сети."
	v.Vars["meta_keywords"] = "технологии, it, стартапы, проекты, новости, форум, сообщество, обсуждения"

	v.Vars["count"] = userCount
	v.Vars["adminCount"] = adminCount

	return v.Render(context)
}
