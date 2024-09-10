package inspections

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/rmishgoog/adopt-a-dog/app/api/errs"
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

func genErrors(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return errs.Newf(errs.FailedPrecondition, "this message is trused")
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, r, status, http.StatusOK)
}

func genPanics(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	return web.Respond(ctx, w, r, status, http.StatusOK)
}
