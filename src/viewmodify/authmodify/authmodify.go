package authmodify

import (
	"github.com/techfront/core/src/kernel/view"
)

/**
 * Модификатор шаблонов, получает токен из контекста и передает переменную в шаблонизатор.
 */
func Modify(context view.Context, v *view.View) {
	v.Vars["authenticity_token"] = context.Get("authenticity_token")
}
