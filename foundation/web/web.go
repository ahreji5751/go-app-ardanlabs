package web

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type Logger func(ctx context.Context, msg string, v ...any)

type App struct {
	*http.ServeMux
	shutdown chan os.Signal
	mv       []MidHandler
	log      Logger
}

func NewApp(log Logger, shutdown chan os.Signal, mv ...MidHandler) *App {
	return &App{
		log:      log,
		ServeMux: http.NewServeMux(),
		shutdown: shutdown,
		mv:       mv,
	}
}

func (a *App) HandleFunc(pattern string, withMiddleware bool, handler Handler, mv ...MidHandler) {
	if withMiddleware {
		handler = wrapMiddleware(mv, handler)
		handler = wrapMiddleware(a.mv, handler)
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}
		ctx := setValues(r.Context(), &v)

		if err := handler(ctx, w, r); err != nil {
			a.log(ctx, "web", "ERROR", err)
			return
		}
	}

	a.ServeMux.HandleFunc(pattern, h)
}
