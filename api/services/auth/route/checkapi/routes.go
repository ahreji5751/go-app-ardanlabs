package checkapi

import (
	"github.com/ardanlabs/service/app/api/auth"
	"github.com/ardanlabs/service/foundation/web"
)

func Routes(app *web.App, ath *auth.Auth) {
	app.HandleFunc("GET /liveness", false, liveness)
	app.HandleFunc("GET /readiness", false, readiness)
}
