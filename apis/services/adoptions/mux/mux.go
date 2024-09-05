package mux

import (
	"os"

	"github.com/rmishgoog/adopt-a-dog/apis/services/adoptions/route/system/healthchek"
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

func WebAPI(shutdown chan os.Signal) *web.App {

	// It's always a good practice to create your own mux, never use the DefaultServerMux for production projects.
	mux := web.NewApp(shutdown)
	healthchek.Routes(mux)

	return mux
}
