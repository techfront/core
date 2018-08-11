package adminaction

import (
	"fmt"
	"time"

	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"

	"github.com/techfront/core/src/component/topic"
)

const DEFAULT_LIMIT = 12

func HandleTopic(context router.Context) error {

	// Проверка прав
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err, displayerror.AccessDenied...)
	}

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Построение запроса
	query := topic.Query().Limit(DEFAULT_LIMIT)

	query.Order("topic_id desc")

	filter := params.Get("status")
	if len(filter) > 0 {
		query.Where("topic_status=?", filter)
	}

	// Определяем текущую страницу и делаем отступ
	currentPage := int(params.GetInt("page"))
	if currentPage > 0 {
		query.Offset(DEFAULT_LIMIT * currentPage)
	}

	// Получение запроса
	results, err := topic.FindAll(query)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Отображение шаблона
	v := view.New("component/admin/template/topics/index", "component/admin/template/topics/row_mini")

	v.Vars["page"] = currentPage
	v.Vars["topics"] = results
	v.Vars["topics_count"] = len(results)

	if err := setTopicMetadata(v, context); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return v.Render(context)
}

/**
* Функция задаёт базовые параметры для топиков.
*
* @param v *view.View Содержит инициализированый View.
* @param context router.Context Содержит интерфейс контекста.
 */
func setTopicMetadata(v *view.View, context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return err
	}

	// Определение следующей страницы
	nextPage := params.GetInt("page") + 1
	prevPage := params.GetInt("page") - 1

	// Задаем базовые параметры для топиков
	v.Vars["pub_date"] = time.Now()
	v.Vars["next_page"] = nextPage
	v.Vars["prev_page"] = prevPage

	v.Vars["next_page_link"] = fmt.Sprintf("?page=%d", nextPage)
	v.Vars["prev_page_link"] = fmt.Sprintf("?page=%d", prevPage)

	filter := params.Get("status")
	if len(filter) > 0 {
		v.Vars["next_page_link"] = fmt.Sprintf("?status=%s&page=%d", filter, nextPage)
		v.Vars["prev_page_link"] = fmt.Sprintf("?status=%s&page=%d", filter, prevPage)
	}
	
	v.Vars["meta_title"] = "Топики / Техфронт"
	v.Vars["meta_desc"] = "Techfront - сообщество тех, кто любит обсуждать технологии, проекты и всё то, что происходит в интернете."
	v.Vars["meta_keywords"] = "технологии, it, стартапы, проекты, форум, сообщество, обсуждения"

	return nil
}
