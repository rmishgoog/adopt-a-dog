package middleware

import (
	"context"

	"github.com/rmishgoog/adopt-a-dog/app/api/metrics"
)

func Metrics(ctx context.Context, handler Handler) error {

	ctx = metrics.Set(ctx)
	err := handler(ctx)
	n := metrics.AddRequest(ctx)
	if n%100 == 0 {
		metrics.AddGoRoutines(ctx)
	}

	if err != nil {
		metrics.AddError(ctx)
	}

	return err
}
