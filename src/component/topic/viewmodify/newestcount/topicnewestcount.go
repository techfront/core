package newestcount

import (
	"github.com/techfront/core/src/component/topic"
	"github.com/techfront/core/src/kernel/view"
)

/**
* Модификатор шаблонов, создает переменную с колличеством новых топиков за промежуток времени.
 */
func Modify(context view.Context, v *view.View) {
	v.Vars["topic_newest_count"] = topic.GetNewestCount("86400s")
}
