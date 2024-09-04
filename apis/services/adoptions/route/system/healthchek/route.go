package healthchek

import "net/http"

// Build routes and configure them on the mux.
func Routes(mux *http.ServeMux) {

	mux.HandleFunc("/liveness", liveness)
	mux.HandleFunc("/readiness", readiness)
}
