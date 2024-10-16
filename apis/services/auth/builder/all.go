package builder

import (
	"github.com/rmishgoog/adopt-a-dog/apis/services/adoptions/route/system/healthchek"
	"github.com/rmishgoog/adopt-a-dog/apis/services/auth/mux"
	"github.com/rmishgoog/adopt-a-dog/apis/services/auth/route/domain/auth"
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

func Routes() adder {
	return adder{}
}

type adder struct {
}

func (adder adder) Add(app *web.App, cfg mux.Config) {

	healthchek.Routes(app, healthchek.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
	})

	auth.Routes(app, auth.Config{
		Auth: cfg.Auth,
	})
}
