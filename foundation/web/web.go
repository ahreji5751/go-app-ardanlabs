package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	*http.ServeMux
	shutdown chan os.Signal
	mv       []MidHandler
}

func NewApp(shutdown chan os.Signal, mv ...MidHandler) *App {
	return &App{
		ServeMux: http.NewServeMux(),
		shutdown: shutdown,
		mv:       mv,
	}
}

func (a *App) HandleFunc(pattern string, handler Handler, mv ...MidHandler) {
	handler = wrapMiddleware(mv, handler)
	handler = wrapMiddleware(a.mv, handler)
	h := func(w http.ResponseWriter, r *http.Request) {
		if err := handler(r.Context(), w, r); err != nil {
			fmt.Println(err)
			return
		}
	}

	a.ServeMux.HandleFunc(pattern, h)
}
