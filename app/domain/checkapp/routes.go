package chekapp

import (
	"github.com/ardanlabs/service/foundation/web"
)

func Routes(app *web.App) {
	app.HandleFunc("GET /liveness", liveness)
	app.HandleFunc("GET /readiness", readiness)
}
