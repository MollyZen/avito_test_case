package v1

import (
	"avito_test_case/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ironstar-io/chizerolog"
)

func NewRouter(l logger.Logger) *chi.Mux {
	r := chi.NewRouter()

	switch v := l.(type) {
	case *logger.ZeroLogLogger:
		r.Use(chizerolog.LoggerMiddleware(v.L))
	default:
		r.Use(middleware.Logger)
	}

	r.Use(middleware.Recoverer)

	return r
}
