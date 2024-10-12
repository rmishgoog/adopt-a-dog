package mux

import (
	"os"

	coreauth "github.com/rmishgoog/adopt-a-dog/core/api/auth"
	"github.com/rmishgoog/adopt-a-dog/foundations/logger"
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

type Config struct {
	Build string
	Log   *logger.Logger
	Auth  *coreauth.Auth
}

func WebAPI(cfg Config, shutdown chan os.Signal) *web.App {

	//mux := web.NewApp(shutdown, middleware.Logger(cfg.Log), middleware.Errors(cfg.Log), middleware.Metrics(), middleware.Panics())

	//healthchek.Routes(mux)
	//auth.Routes(mux, auth.Config{Auth: cfg.Auth})

	//return mux
	return nil
}
