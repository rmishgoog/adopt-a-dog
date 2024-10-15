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
	// Passing in a MidHandler, which takes a real handler & returns another handler after wrapping it w/
	// the middleware logic it needs to execute, in this case, it's the authentication.
	// Here the real handler is api.profile which will be wrapped w/ invocation of token validation.
	app.HandleFunc("GET /profile/{uid}", api.profile, authnmid)
}
