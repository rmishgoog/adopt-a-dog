package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	*http.ServeMux
	shutdown chan os.Signal
	mw       []MidHandler
}

func NewApp(shutdown chan os.Signal, mw ...MidHandler) *App {
	return &App{
		ServeMux: http.NewServeMux(),
		shutdown: shutdown,
		mw:       mw,
	}
}

// This is an override of the promoted API method from the embedded ServeMux.
func (a *App) HandleFunc(pattern string, handler Handler, mw ...MidHandler) {

	handler = wrapMiddleware(mw, handler)   // At this point in time we have none, so this is a no-op.
	handler = wrapMiddleware(a.mw, handler) // Here we have the middleware from the App struct, currently the logger middleware.

	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}
		ctx := setValues(r.Context(), &v)
		if err := handler(ctx, w, r); err != nil {
			fmt.Println(err)
		}
	}
	a.ServeMux.HandleFunc(pattern, h)
}

func (a *App) HandleFuncNoMiddleware(pattern string, handler Handler, mw ...MidHandler) {

	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}
		ctx := setValues(r.Context(), &v)
		if err := handler(ctx, w, r); err != nil {
			fmt.Println(err)
		}
	}
	a.ServeMux.HandleFunc(pattern, h)
}
