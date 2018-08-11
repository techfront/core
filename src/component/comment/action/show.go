package commentaction

import (
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/component/comment"

	"github.com/techfront/core/src/lib/displayerror"
)

/**
* HandleShow displays a single comment
 */
func HandleShow(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Получение комментария
	commentEntity, err := comment.Find(params.GetInt("id"))
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	// Отображение шаблона
	v := view.New("component/comment/template/show")

	v.Vars["comment"] = commentEntity

	return v.Render(context)
}
