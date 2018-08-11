package msgmodify

import (
	"github.com/techfront/core/src/kernel/view"
)

/**
 * Модификатор шаблонов, получает сообщение из контекста и передает переменные в шаблонизатор.
 */
func Modify(context view.Context, v *view.View) {
	m := context.Message()

	if v.Vars["message"] == nil {
		v.Vars["message"] = m.Body()

		var s string
		switch m.Style() {
		case 0:
			s = "default"
		case 1:
			s = "error"
		case 2:
			s = "success"
		case 3:
			s = "warning"
		default:
			s = "default"
		}

		v.Vars["message_class"] = "js-message--" + m.Uid()
		v.Vars["message_type"] = s
	}
}
