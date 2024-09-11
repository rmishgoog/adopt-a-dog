package metrics

import (
	"context"
	"expvar"
	"runtime"
)

var m metrics

type metrics struct {
	goroutines *expvar.Int
	errors     *expvar.Int
	panics     *expvar.Int
	requests   *expvar.Int
}

func init() {
	m = metrics{
		goroutines: expvar.NewInt("goroutines"),
		errors:     expvar.NewInt("errors"),
		panics:     expvar.NewInt("panics"),
		requests:   expvar.NewInt("requests"),
	}
}

type ctxKey int

const key ctxKey = 1

func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, &m)
}

func AddGoRoutines(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		g := int64(runtime.NumGoroutine())
		v.goroutines.Set(g)
		return g
	}
	return 0
}

func AddError(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.errors.Add(1)
		return v.errors.Value()
	}
	return 0
}

func AddPanic(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.panics.Add(1)
		return v.panics.Value()
	}
	return 0
}

func AddRequest(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.requests.Add(1)
		return v.requests.Value()
	}
	return 0
}
