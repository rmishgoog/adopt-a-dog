package profile

import "github.com/rmishgoog/adopt-a-dog/foundations/web"

func Routes(app *web.App) {

	app.HandleFunc("/profile/{user_id}", profile)
}
