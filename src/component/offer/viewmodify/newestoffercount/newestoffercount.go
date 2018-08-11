package newestoffercount

import (
	"github.com/techfront/core/src/component/offer"
	"github.com/techfront/core/src/kernel/view"
)

/**
* Модификатор шаблонов, создает переменную с колличеством новых офферов за промежуток времени.
 */
func Modify(context view.Context, v *view.View) {
	v.Vars["offer_newest_count"] = offer.GetNewestCount("86400s")
}
