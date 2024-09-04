package mux

import (
	"net/http"

	"github.com/rmishgoog/adopt-a-dog/apis/services/adoptions/route/system/healthchek"
)

func WebAPI() *http.ServeMux {

	// It's always a good practice to create your own mux, never use the DefaultServerMux for production projects.
	mux := http.NewServeMux()
	healthchek.Routes(mux)

	return mux
}
