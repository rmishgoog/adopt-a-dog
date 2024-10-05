package auth

import (
	"github.com/rmishgoog/adopt-a-dog/core/api/auth"
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

type Config struct {
	Auth *auth.Auth
}

func Routes(app *web.App, cfg Config) {

	api := newAPI(cfg.Auth)
	//In Go 1.22 and later it is allowed to specify the method in the HandleFunc call.
	app.HandleFunc("POST /authenticate", api.authenticate)
}
