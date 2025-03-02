package mux

import (
	"github.com/ardanlabs/service/app/domain/checkapp"
	"github.com/ardanlabs/service/foundation/web"
	"os"
)

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(shutdown chan os.Signal) *web.App {
	mux := web.NewApp(shutdown)

	chekapp.Routes(mux)

	return mux
}
