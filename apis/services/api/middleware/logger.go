package middleware

import (
	"context"
	"net/http"

	"github.com/rmishgoog/adopt-a-dog/app/api/middleware"
	"github.com/rmishgoog/adopt-a-dog/foundations/logger"
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

func Logger(log *logger.Logger) web.MidHandler {

	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// Uses Go closures here.
			handle := func(ctx context.Context) error {
				// Read this like:
				// return liveness(ctx, w, r) or readiness (ctx, w, r) ~~ baked in here.
				return handler(ctx, w, r)
			}
			// This is the actual middleware and upon a h(), this statement is executed, invoking the closure based baked-in code above.
			return middleware.Logger(ctx, log, r.URL.Path, r.URL.RawQuery, r.Method, r.RemoteAddr, handle)
		}
		return h
	}
	return m

}
