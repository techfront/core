package pageaction

import (
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"
)

func HandleSpasibo(context router.Context) error {

	v := view.New("component/page/template/spasibo")

	r := context.Request()

	v.Vars["meta_title"] = "Спасибо / Техфронт"
	v.Vars["meta_desc"] = "Сообщество энтузиастов, кому по нраву обсуждать технологии, исследования, стартапы и всё то, что происходит в сети."
	v.Vars["meta_keywords"] = "технологии, it, стартапы, проекты, новости, форум, сообщество, обсуждения"

	v.Vars["meta_refresh"] = r.URL.Scheme + r.Host
	v.Vars["meta_refresh_s"] = "2"

	return v.Render(context)
}
