package checkapi

import (
	"github.com/ardanlabs/service/foundation/web"
)

func Routes(app *web.App) {
	app.HandleFunc("GET /liveness", false, liveness)
	app.HandleFunc("GET /readiness", false, readiness)
	app.HandleFunc("GET /test-error", true, testError)
	app.HandleFunc("GET /test-panic", true, testPanic)
}
