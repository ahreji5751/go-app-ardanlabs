package checkapi

import (
	"github.com/ardanlabs/service/api/services/api/mid"
	"github.com/ardanlabs/service/app/api/auth"
	"github.com/ardanlabs/service/app/api/authclient"
	"github.com/ardanlabs/service/foundation/logger"
	"github.com/ardanlabs/service/foundation/web"
)

func Routes(app *web.App, log *logger.Logger, authClient *authclient.Client) {
	authenticate := mid.Authenticate(log, authClient)
	authorizeAdmin := mid.Authorize(log, authClient, auth.RuleAdminOnly)

	app.HandleFunc("GET /liveness", false, liveness)
	app.HandleFunc("GET /readiness", false, readiness)
	app.HandleFunc("GET /test-error", true, testError)
	app.HandleFunc("GET /test-panic", true, testPanic)
	app.HandleFunc("GET /test-auth", true, liveness, authenticate, authorizeAdmin)
}
