package currentpath

import (
	"github.com/techfront/core/src/kernel/view"
)

/**
 * Модификатор шаблонов, создает переменную с текущем путем.
 * Для определения активной страницы.
 */
func Modify(context view.Context, v *view.View) {
	v.Vars["current_path"] = context.Path()
}
