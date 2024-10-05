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
	app.HandleFunc("POST /authenticate", api.authenticate)
}
