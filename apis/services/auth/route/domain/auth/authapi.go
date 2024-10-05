package auth

import (
	"context"
	"net/http"

	"github.com/rmishgoog/adopt-a-dog/core/api/auth"
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

type api struct {
	auth *auth.Auth
}

func newAPI(auth *auth.Auth) *api {
	return &api{
		auth: auth,
	}
}

func (a *api) authenticate(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	autheHeader := r.Header.Get("Authorization")
	if autheHeader == "" {
		status := struct {
			Status string `json:"status"`
		}{
			Status: "unauthorized",
		}
		return web.Respond(ctx, w, r, status, http.StatusUnauthorized)
	}

	return nil
}
