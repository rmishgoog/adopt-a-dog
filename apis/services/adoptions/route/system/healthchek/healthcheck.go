package healthchek

import (
	"context"
	"net/http"

	"github.com/rmishgoog/adopt-a-dog/foundations/logger"
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

type api struct {
	build string
	log   *logger.Logger
}

func newAPI(build string, log *logger.Logger) *api {
	return &api{
		build: build,
		log:   log,
	}
}

func (api *api) liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	return web.Respond(ctx, w, r, status, http.StatusOK)
}

func (api *api) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	return web.Respond(ctx, w, r, status, http.StatusOK)
}
