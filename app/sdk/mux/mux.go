package mux

import (
	"github.com/ardanlabs/service/app/domain/checkapp"
	"net/http"
)

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI() *http.ServeMux {
	mux := http.NewServeMux()

	chekapp.Routes(mux)

	return mux
}
