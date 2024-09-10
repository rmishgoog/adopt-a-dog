package middleware

import (
	"context"
	"fmt"
	"runtime/debug"
)

func Panics(ctx context.Context, handler Handler) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			trace := debug.Stack()
			err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))
		}
	}()

	return handler(ctx)
}
