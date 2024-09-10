package middleware

import (
	"context"
	"net/http"

	"github.com/rmishgoog/adopt-a-dog/app/api/middleware"
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

func Panics() web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return middleware.Panics(ctx, hdl)
		}

		return h
	}

	return m
}
