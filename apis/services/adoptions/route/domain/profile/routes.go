package profile

import (
	"github.com/rmishgoog/adopt-a-dog/apis/services/api/middleware"
	"github.com/rmishgoog/adopt-a-dog/app/api/authclient"
	"github.com/rmishgoog/adopt-a-dog/foundations/logger"
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

type Config struct {
	Build      string
	Log        *logger.Logger
	AuthClient *authclient.ServiceClient
}

func Routes(app *web.App, cfg Config) {

	authnmid := middleware.Authenticate(cfg.Log, cfg.AuthClient)

	api := newAPI(cfg.Build, cfg.Log, cfg.AuthClient)
	app.HandleFunc("GET /profile/{uid}", api.profile, authnmid)
}
