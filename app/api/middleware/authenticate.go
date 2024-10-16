package middleware

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rmishgoog/adopt-a-dog/app/api/authclient"
	"github.com/rmishgoog/adopt-a-dog/app/api/errs"
	"github.com/rmishgoog/adopt-a-dog/core/api/auth"
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

// Middleware function at the app (protocol agnostic) to process the JWT token issued by keycloak server.
func Bearer(ctx context.Context, auth *auth.Auth, authorization string, handler Handler) error {

	claims, err := auth.Authenticate(ctx, authorization)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	if claims.Subject == "" {
		return errs.Newf(errs.Unauthenticated, "absent claims: you are not authorized for that action, no claims")
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return errs.New(errs.Unauthenticated, fmt.Errorf("parsing subject failed: %w", err))
	}

	ctx = setUserID(ctx, subjectID)
	ctx = setClaims(ctx, claims)

	return handler(ctx)
}
