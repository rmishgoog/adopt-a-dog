package mux

import (
	"os"

	"github.com/rmishgoog/adopt-a-dog/apis/services/api/middleware"
	healthchek "github.com/rmishgoog/adopt-a-dog/apis/services/auth/route/system/healthcheck"
	"github.com/rmishgoog/adopt-a-dog/core/api/auth"
	"github.com/rmishgoog/adopt-a-dog/foundations/logger"
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

func WebAPI(log *logger.Logger, auth *auth.Auth, shutdown chan os.Signal) *web.App {

	mux := web.NewApp(shutdown, middleware.Logger(log), middleware.Errors(log), middleware.Metrics(), middleware.Panics())
	healthchek.Routes(mux)

	return mux
}
