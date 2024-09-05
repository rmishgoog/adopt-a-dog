package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	*http.ServeMux
	shutdown chan os.Signal
}

func NewApp(shutdown chan os.Signal) *App {
	return &App{
		ServeMux: http.NewServeMux(),
		shutdown: shutdown,
	}
}

// This is an override of the promoted API method from the embedded ServeMux.
func (a *App) HandleFunc(path string, handler Handler) {

	h := func(w http.ResponseWriter, r *http.Request) {
		if err := handler(r.Context(), w, r); err != nil {
			fmt.Println(err)
		}
	}
	a.ServeMux.HandleFunc(path, h)
}
