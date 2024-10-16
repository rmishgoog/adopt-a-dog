package auth

import (
	"context"
	"net/http"

	"github.com/rmishgoog/adopt-a-dog/app/api/authclient"
	"github.com/rmishgoog/adopt-a-dog/app/api/errs"
	"github.com/rmishgoog/adopt-a-dog/app/api/middleware"
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

	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	resp := authclient.AuthenticateResp{
		UserID: userID,
		Claims: middleware.GetClaims(ctx),
	}

	return web.Respond(ctx, w, resp, http.StatusOK)
}
