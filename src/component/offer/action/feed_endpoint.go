package offeraction

import (
	"fmt"
	"strings"
	"github.com/techfront/core/src/component/offer"
	"github.com/techfront/core/src/kernel/response"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleFeedEndpoint возвращает ленту офферов в JSON
 */
func HandleFeedEndpoint(context router.Context) error {

	// Получение параметров
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Построение запроса
	query := offer.Query()

	// Обработка limit
	limit := int(params.GetInt("limit"))
	if limit == 0 {
		limit = DEFAULT_LIMIT
	} else if limit > 100 {
		limit = 100
	}

	query.Limit(limit)

	// Обработка page
	page := int(params.GetInt("page"))
	if page > 0 {
		query.Offset(limit * page)
	}

	// Обработка offset
	offset := int(params.GetInt("offset"))
	if offset > 0 {
		query.Offset(limit + offset)
	}

	// Обработка order
	order := params.Get("order")
	if order != "desc" && order != "asc" {
		order = DEFAULT_ORDER
	}
	
	// Обработка sort
	sort := params.Get("sort")
	switch sort {
		case "id":
			sort = fmt.Sprintf("offer_id %s", order)
			query.Order(sort)
		case "score":
			sort = fmt.Sprintf("score(offer_count_upvote, offer_count_downvote, offer_count_flag) %s", order)
			query.Order(sort)
		case "rank":
			sort = fmt.Sprintf("rank(score(offer_count_upvote, offer_count_downvote, offer_count_flag), offer_created_at) %s", order)
			query.Order(sort)
		default:
			sort = fmt.Sprintf("offer_id %s", order)
			query.Order(sort)
	}

	// Обработка q
	filter := params.Get("q")
	if len(filter) > 0 {
		filter = strings.Replace(filter, "_", "\\_", -1)
		filter = strings.Replace(filter, "%", "\\%", -1)
		wildcard := "%" + filter+ "%"

		query.Where("offer_name ILIKE ? OR offer_url ILIKE ?", wildcard, wildcard)
	}

	// Обработка status
	status := params.Get("status")
	if len(status) > 0 {
		var statusArr []string
		for _, v := range strings.Split(status, ",") {
			switch v {
			case "100":
				statusArr = append(statusArr, "offer_status = 100")
			case "101":
				statusArr = append(statusArr, "offer_status = 101")
			case "13":
				statusArr = append(statusArr, "offer_status = 13" )
			case "15":
				statusArr = append(statusArr, "offer_status = 15")
			}
		}

		statusWhere := strings.Join(statusArr, " OR ")
		query.Where(statusWhere)
		
	}

	// Обработка format
	format := params.Get("format")
	switch format {
		case "offer":
			query.Where("offer_id_format = 0")
		case "news":
			query.Where("offer_id_format = 10")
		case "video":
			query.Where("offer_id_format = 20")
		case "question":
			query.Where("offer_id_format = 30")
		case "project":
			query.Where("offer_id_format = 40")
		case "podcast":
			query.Where("offer_id_format = 50")
	}

	// Получение офферов
	results, err := offer.FindAll(query)
	if err != nil {
		return response.Send(context.Writer(), 404, response.NotFound("Error 404"))
	}

	return response.Send(context.Writer(), 200, &response.Response{
		Data: results,
		Meta: map[string]interface{}{},
	})
}
