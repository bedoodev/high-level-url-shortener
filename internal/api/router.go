package api

import (
	"net/http"

	_ "github.com/bedoodev/high-level-url-shortener/docs"
	"github.com/bedoodev/high-level-url-shortener/internal/middleware"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(handler *Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.NewRateLimiterMiddleware("30-M"))

	r.Get("/ping", handler.Ping)

	r.Post("/shorten", handler.Shorten)
	r.Get("/{code}", handler.Resolve)
	r.Get("/analytics/{shortCode}", handler.GetAnalytics)
	r.Get("/analytics/top", handler.GetTopAnalytics)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	return r
}
