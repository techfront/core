package redirecthandler

import (
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/displayerror"
	"net/url"
)

func Handler(context router.Context) error {
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	to := params.Get("to")

	// Проверка URL
	if _, err := url.Parse(to); err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	return router.RedirectExternal(context, to)
}
