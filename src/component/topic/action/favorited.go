package topicaction

import (
	"strings"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"
	"github.com/techfront/core/src/component/topic"
	"github.com/techfront/core/src/lib/authorise"
	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleFavorites отображает список топиков из закладок.
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
	query := topic.Favorited(userEntity.Id).Limit(DEFAULT_LIMIT)

	// Фильтр по q (поисковой запрос)
	qFilter := params.Get("q")
	if len(qFilter) > 0 {
		qFilter = strings.Replace(qFilter, "_", "\\_", -1)
		qFilter = strings.Replace(qFilter, "%", "\\%", -1)
		wildcard := "%" + qFilter + "%"

		query.Where("topic_name ILIKE ? OR topic_url ILIKE ?", wildcard, wildcard)
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
	v := view.New(
		"component/topic/template/index",
		"component/topic/template/row_mini",
	)

	v.Vars["page"] = currentPage
	v.Vars["topics"] = results
	v.Vars["topics_format"] = -1
	v.Vars["topics_count"] = len(results)

	if err := setTopicMetadata(v, context); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	v.Vars["page_icon"] = "icon-bookmark"
	v.Vars["page_title"] = "Избранное"
	v.Vars["meta_title"] = "Избранное / Техфронт"

	return v.Render(context)
}