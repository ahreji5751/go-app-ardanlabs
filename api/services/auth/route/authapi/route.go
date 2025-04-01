package authapi

import (
	"github.com/ardanlabs/service/api/services/api/mid"
	"github.com/ardanlabs/service/app/api/auth"
	"github.com/ardanlabs/service/foundation/web"
)

func Routes(app *web.App, ath *auth.Auth) {
	bearer := mid.Bearer(ath)
	basic := mid.Basic(ath)

	api := newAPI(ath)
	app.HandleFunc("GET /auth/token/{kid}", true, api.token, basic)
	app.HandleFunc("GET /auth/authenticate", true, api.authenticate, bearer)
	app.HandleFunc("POST /auth/authorize", true, api.authorize)
}
