package pageaction

import (
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"
)

func HandleWaysToGive(context router.Context) error {

	v := view.New("component/page/template/ways_to_give")

	v.Vars["meta_title"] = "Способы поддержать проект / Техфронт"
	v.Vars["meta_desc"] = "Сообщество энтузиастов, кому по нраву обсуждать технологии, исследования, стартапы и всё то, что происходит в сети."
	v.Vars["meta_keywords"] = "технологии, it, стартапы, проекты, новости, форум, сообщество, обсуждения"

	return v.Render(context)
}
