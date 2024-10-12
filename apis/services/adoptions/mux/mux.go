package mux

import (
	"context"
	"os"

	"github.com/rmishgoog/adopt-a-dog/apis/services/api/middleware"
	"github.com/rmishgoog/adopt-a-dog/app/api/authclient"

	"github.com/rmishgoog/adopt-a-dog/foundations/logger"
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

type Config struct {
	Build      string
	Log        *logger.Logger
	AuthClient *authclient.ServiceClient
	Shutdown   chan os.Signal
}

type RouteAdder interface {
	Add(app *web.App, cfg Config)
}

func WebAPI(cfg Config, adder RouteAdder) *web.App {

	logger := func(ctx context.Context, msg string, v ...any) {
		cfg.Log.Info(ctx, msg, v...)
	}

	mux := web.NewApp(cfg.Shutdown,
		logger, middleware.Logger(cfg.Log),
		middleware.Errors(cfg.Log),
		middleware.Metrics(),
		middleware.Panics(),
	)
	adder.Add(mux, cfg)

	return mux
}
