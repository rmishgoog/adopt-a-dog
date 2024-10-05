package healthchek

import (
	"context"
	"net/http"

	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

func liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	return web.Respond(ctx, w, r, status, http.StatusOK)
}

func readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	return web.Respond(ctx, w, r, status, http.StatusOK)
}
