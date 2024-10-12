package profile

import (
	"context"
	"net/http"

	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

// This handler function will load an existing profile for an end-user
func profile(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	return web.Respond(ctx, w, r, status, http.StatusOK)
}
