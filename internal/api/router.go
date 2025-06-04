package api

import (
	"net/http"

	_ "github.com/bedoodev/high-level-url-shortener/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(handler *Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/ping", handler.Ping)

	r.Post("/shorten", handler.Shorten)
	r.Get("/{code}", handler.Resolve)

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	return r
}
