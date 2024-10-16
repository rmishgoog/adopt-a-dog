package profile

import (
	"context"
	"net/http"

	"github.com/rmishgoog/adopt-a-dog/app/api/authclient"
	"github.com/rmishgoog/adopt-a-dog/foundations/logger"
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

type api struct {
	build      string
	log        *logger.Logger
	authClient *authclient.ServiceClient
}

func newAPI(build string, log *logger.Logger, authclient *authclient.ServiceClient) *api {
	return &api{
		build:      build,
		log:        log,
		authClient: authclient,
	}
}

// TODO This handler function will load an existing profile for an end-user
func (api *api) profile(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	return web.Respond(ctx, w, status, http.StatusOK)
}
