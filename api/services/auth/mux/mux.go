package mux

import (
	"context"
	"os"

	"github.com/ardanlabs/service/api/services/api/mid"
	"github.com/ardanlabs/service/api/services/auth/route/authapi"
	"github.com/ardanlabs/service/api/services/auth/route/checkapi"
	"github.com/ardanlabs/service/app/api/auth"
	"github.com/ardanlabs/service/foundation/logger"
	"github.com/ardanlabs/service/foundation/web"
)

func WebAPI(log *logger.Logger, ath *auth.Auth, shutdown chan os.Signal) *web.App {
	loggerFn := func(ctx context.Context, msg string, v ...any) {
		log.Info(ctx, msg, v...)
	}
	app := web.NewApp(loggerFn, shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics())

	checkapi.Routes(app, ath)
	authapi.Routes(app, ath)

	return app
}
