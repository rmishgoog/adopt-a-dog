package healthchek

import (
	"github.com/rmishgoog/adopt-a-dog/foundations/logger"
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

type Config struct {
	Build string
	Log   *logger.Logger
}

// Build routes and configure them on the mux.
func Routes(app *web.App, cfg Config) {

	api := newAPI(cfg.Build, cfg.Log)

	app.HandleFuncNoMiddleware("/liveness", api.liveness)
	app.HandleFuncNoMiddleware("/readiness", api.readiness)
}
