package errorhandler

import (
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view"
)

func Handler(context router.Context, e error) {
	context.Logf("%s", e)

	err := router.ToStatusError(e)

	context.Logf("#error Status %d at %s", err.Status, err.FileLine())

	v := view.New("template/error")

	v.Templates = []string{"template/meta", "template/includes", "template/error"}

	v.Layout = "error"

	v.Vars["error_title"] = err.Title
	v.Vars["error_message"] = err.Message

	if !context.Production() {
		v.Vars["error_status"] = err.Status
		v.Vars["error_file"] = err.FileLine()
		v.Vars["error_error"] = err.Err
	}

	v.Render(context)
}
