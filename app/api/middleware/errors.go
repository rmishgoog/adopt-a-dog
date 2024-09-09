package middleware

import (
	"context"

	"github.com/rmishgoog/adopt-a-dog/app/api/errs"
	"github.com/rmishgoog/adopt-a-dog/foundations/logger"
)

func Errors(ctx context.Context, log *logger.Logger, handler Handler) error {

	err := handler(ctx) // This is important in the app layer middleware, as this is where you hide the protocol specific details.
	if err == nil {
		return nil
	}
	log.Error(ctx, "message", "ERROR", err.Error())
	if errs.IsError(err) {
		return errs.GetError(err)
	}
	return errs.Newf(errs.Unknown, errs.Unknown.String())
}
