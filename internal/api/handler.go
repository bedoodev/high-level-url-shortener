package api

import (
	"encoding/json"
	"net/http"

	"github.com/bedoodev/high-level-url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	urlService service.URLService
}

func NewHandler(service service.URLService) *Handler {
	return &Handler{urlService: service}
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

// Shorten godoc
// @Summary     Shorten a URL
// @Description Takes a long URL and returns a shortened version
// @Accept      json
// @Produce     json
// @Param       request body shortenRequest true "URL to shorten"
// @Success     200 {object} shortenResponse
// @Failure     400 {string} string "Invalid request"
// @Failure     500 {string} string "Internal server error"
// @Router      /shorten [post]
func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.urlService.ShortenURL(r.Context(), req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := shortenResponse{
		ShortURL: "http://localhost:8080/" + result.ShortCode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Resolve godoc
// @Summary     Redirect to original URL
// @Description Redirects from a short URL code to the original URL
// @Produce     plain
// @Param       code path string true "Short code"
// @Success     302 {string} string "Redirect"
// @Failure     404 {string} string "Not found"
// @Router      /{code} [get]
func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	result, err := h.urlService.ResolveURL(r.Context(), code)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, result.OriginalURL, http.StatusFound)
}

// HealthCheck godoc
// @Summary     HealthCheck
// @Success     200 {string} string "OK"
// @Failure     500 {string} string "Internal Server Error"
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
