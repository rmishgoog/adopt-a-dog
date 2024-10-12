package builder

import (
	"github.com/rmishgoog/adopt-a-dog/apis/services/adoptions/mux"
	"github.com/rmishgoog/adopt-a-dog/apis/services/adoptions/route/domain/profile"
	"github.com/rmishgoog/adopt-a-dog/apis/services/adoptions/route/system/healthchek"
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

	profile.Routes(app, profile.Config{
		Build:      cfg.Build,
		Log:        cfg.Log,
		AuthClient: cfg.AuthClient,
	})

}
