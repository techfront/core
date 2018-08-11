package usermodify

import (
	"github.com/techfront/core/src/kernel/view"
)

/**
 * Модификатор шаблонов, получает объект текущего пользователя из контекста и передает переменную в шаблонизатор.
 */
func Modify(context view.Context, v *view.View) {
	v.Vars["current_user"] = context.Get("current_user")
}
