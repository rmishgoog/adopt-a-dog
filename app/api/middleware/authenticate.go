package middleware

import (
	"context"

	"github.com/rmishgoog/adopt-a-dog/app/api/authclient"
	"github.com/rmishgoog/adopt-a-dog/app/api/errs"
	"github.com/rmishgoog/adopt-a-dog/foundations/logger"
)

func Authenticate(ctx context.Context, log *logger.Logger, client *authclient.ServiceClient, authorization string, hdl Handler) error {

	resp, err := client.Authenticate(ctx, authorization)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	ctx = setUserID(ctx, resp.UserID)
	ctx = setClaims(ctx, resp.Claims)

	return hdl(ctx)
}
