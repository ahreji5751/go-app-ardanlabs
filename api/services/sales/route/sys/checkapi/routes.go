package checkapi

import (
	"github.com/ardanlabs/service/foundation/web"
)

func Routes(app *web.App) {
	app.HandleFunc("GET /liveness", liveness)
	app.HandleFunc("GET /readiness", readiness)
	app.HandleFunc("GET /test-error", testError)
	app.HandleFunc("GET /test-panic", testPanic)
}
