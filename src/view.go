package src

import (
	"time"
	"github.com/fragmenta/server"
	"github.com/techfront/core/src/kernel/view"

	"github.com/techfront/core/src/viewhelper/displayfloat"
	"github.com/techfront/core/src/viewhelper/fetchquery"
	"github.com/techfront/core/src/viewhelper/isotime"
	"github.com/techfront/core/src/viewhelper/markup"
	"github.com/techfront/core/src/viewhelper/resizeimage"
	"github.com/techfront/core/src/viewhelper/timeago"
	"github.com/techfront/core/src/viewhelper/cleanhtml"

	"github.com/techfront/core/src/component/topic/viewmodify/newestcount"
	"github.com/techfront/core/src/component/offer/viewmodify/newestoffercount"
	"github.com/techfront/core/src/viewmodify/currentpath"
	"github.com/techfront/core/src/viewmodify/msgmodify"
	"github.com/techfront/core/src/viewmodify/authmodify"
	"github.com/techfront/core/src/viewmodify/usermodify"
	"github.com/techfront/core/src/viewmodify/recentcomments"
)

/**
 * Инициализация и конфигурирование шаблонизатора.
 */
func setupView(server *server.Server) {
	defer server.Timef("#info Finished loading templates in %s", time.Now())

	viewConfig := view.View{
		Extension:  "html.got",
		Layout:     "layout",
		Format:     "text/html",
		Folder:     "src",
		Production: server.Production(),
	}

	view.LoadTemplates([]string{
		"template/layout",
		"template/meta",
		"template/includes",
		"template/header",
		"template/footer",
		"template/modal",
		"template/sidebar",

		"template/widget/widget_fixed",
		"template/widget/widget_sidebar_ads",
		"template/widget/widget_sidebar_social",

		"component/user/template/widget/widget_subscribe",
		"component/comment/template/widget/widget_activity",
		"component/comment/template/widget/widget_activity_row",
	})

	view.LoadHelpers(
		timeago.Helper(),
		fetchquery.Helper(),
		isotime.Helper(),
		markup.Helper(),
		displayfloat.Helper(),
		resizeimage.Helper(),
		cleanhtml.Helper(),
		appAssets.StyleLink(),
		appAssets.ScriptLink(),
	)

	view.LoadModifiers(
		currentpath.Modify,
		recentcomments.Modify,
		newestcount.Modify,
		newestoffercount.Modify,
		authmodify.Modify,
		usermodify.Modify,
		msgmodify.Modify,
	)

	view.Setup(viewConfig)
}
