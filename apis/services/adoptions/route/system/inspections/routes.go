package inspections

import "github.com/rmishgoog/adopt-a-dog/foundations/web"

func Routes(app *web.App) {

	// These are strictly for testing purposes & must be prevented via a Kubernetes Network Policy in production.
	app.HandleFunc("/errors", genErrors)
	app.HandleFunc("/panic", genPanics)
}
