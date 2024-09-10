package mux

import (
	"os"

	"github.com/rmishgoog/adopt-a-dog/apis/services/api/middleware"

	"github.com/rmishgoog/adopt-a-dog/apis/services/adoptions/route/system/healthchek"
	"github.com/rmishgoog/adopt-a-dog/foundations/logger"
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

func WebAPI(log *logger.Logger, shutdown chan os.Signal) *web.App {

	// It's always a good practice to create your own mux, never use the DefaultServerMux for production projects.
	mux := web.NewApp(shutdown, middleware.Logger(log), middleware.Errors(log), middleware.Panics())
	healthchek.Routes(mux)

	return mux
}
