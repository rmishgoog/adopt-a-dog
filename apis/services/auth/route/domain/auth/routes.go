package auth

import (
	"github.com/rmishgoog/adopt-a-dog/apis/services/api/middleware"
	"github.com/rmishgoog/adopt-a-dog/core/api/auth"
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

type Config struct {
	Auth *auth.Auth
}

func Routes(app *web.App, cfg Config) {

	// Obtain a bearer middleware & inject it upfront of the authentication endpoint for JWT validations.
	bearer := middleware.Bearer(cfg.Auth)

	api := newAPI(cfg.Auth)
	app.HandleFunc("GET /authenticate", api.authenticate, bearer)
}
