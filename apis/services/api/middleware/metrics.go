package middleware

import (
	"context"
	"net/http"

	"github.com/rmishgoog/adopt-a-dog/app/api/middleware"
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

func Metrics() web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return middleware.Metrics(ctx, hdl)
		}

		return h
	}

	return m
}
