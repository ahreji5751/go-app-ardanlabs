package chekapp

import "net/http"

func Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /liveness", liveness)
	mux.HandleFunc("GET /readiness", readiness)
}
