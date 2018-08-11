package pprofhandler

import (
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/displayerror"
	"net/http/pprof"
)

func Handler(context router.Context) error {

	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	switch params.Get("pprof") {
	case "/cmdline":
		pprof.Cmdline(context.Writer(), context.Request())
	case "/profile":
		pprof.Profile(context.Writer(), context.Request())
	case "/symbol":
		pprof.Symbol(context.Writer(), context.Request())
	default:
		pprof.Index(context.Writer(), context.Request())
	}

	return nil
}
