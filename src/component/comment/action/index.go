package commentaction

import (
	"fmt"

	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/component/comment"
	"github.com/techfront/core/src/component/user"
	"github.com/techfront/core/src/lib/displayerror"
)

const DEFAULT_LIMIT = 50

/**
* HandleIndex displays a list of comments
 */
func HandleIndex(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	var pageTitle string

	pageTitle = "Дискуссии"

	// Построение запроса
	query := comment.Query().Limit(DEFAULT_LIMIT).Order("comment_created_at desc")

	// Фильтр по u (id пользователя)
	uFilter := params.GetInt("u")
	if uFilter > 0 {
		query.Where("comment_id_user=?", uFilter)

		// Меняем название страницы на пользовательское
		userEntity, err := user.Find(uFilter)
		if err != nil {
			return router.InternalError(err, displayerror.PageNotFound...)
		}
		pageTitle = fmt.Sprintf("Все комментарии %s", userEntity.Name)
	}

	// Определяем текущую страницу и делаем отступ
	currentPage := int(params.GetInt("page"))
	if currentPage > 0 {
		query.Offset(DEFAULT_LIMIT * currentPage)
	}

	// Получение комментариев
	results, err := comment.FindAllWithChild(query)
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Отображение шаблона
	v := view.New(
		"component/comment/template/index",
		"component/comment/template/comment",
		"component/comment/template/form",
		"component/comment/template/form_embed",
	)

	nextPage := params.GetInt("page") + 1
	prevPage := params.GetInt("page") - 1

	v.Vars["next_page"] = nextPage
	v.Vars["prev_page"] = prevPage

	v.Vars["next_page_link"] = fmt.Sprintf("?page=%d", nextPage)
	v.Vars["prev_page_link"] = fmt.Sprintf("?page=%d", prevPage)

	v.Vars["comments"] = results

	v.Vars["page"] = currentPage
	v.Vars["page_title"] = pageTitle
	v.Vars["page_icon"] = "icon-chat"

	v.Vars["meta_title"] = pageTitle + " / Техфронт"
	v.Vars["meta_desc"] = "Сообщество энтузиастов, кому по нраву обсуждать технологии, исследования, стартапы и всё то, что происходит в сети."
	v.Vars["meta_keywords"] = "технологии, it, стартапы, проекты, новости, форум, сообщество, обсуждения"

	return v.Render(context)
}
