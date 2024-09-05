package healthchek

import (
	"github.com/rmishgoog/adopt-a-dog/foundations/web"
)

// Build routes and configure them on the mux.
func Routes(app *web.App) {

	app.HandleFunc("/liveness", liveness)
	app.HandleFunc("/readiness", readiness)
}
