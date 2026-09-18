package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dstepan499-dev/url-shortener-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	service *service.ShortenerService
	baseURL string
}

func NewHandler(service *service.ShortenerService, baseURL string) *Handler {
	return &Handler{
		service: service,
		baseURL: baseURL,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/api/shorten", h.Shorten)
	r.Get("/api/analytics/{alias}", h.GetAnalytics)
	r.Get("/{alias}", h.Redirect)

	return r
}

type ShortenRequest struct {
	URL         string `json:"url"`
	CustomAlias string `json:"custom_alias,omitempty"`
}

type ShortenResponse struct {
	Alias    string `json:"alias"`
	ShortURL string `json:"short_url"`
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	alias, err := h.service.ShortenURl(r.Context(), req.URL, req.CustomAlias)
	if err != nil {
		if errors.Is(err, service.ErrInvalidURL) || errors.Is(err, service.ErrEmptyURL) {
			h.respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Failed to shorten url")
		return
	}

	resp := ShortenResponse{
		Alias:    alias,
		ShortURL: fmt.Sprintf("%s/%s", h.baseURL, alias),
	}

	h.respondWithJSON(w, http.StatusCreated, resp)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")
	if alias == "" {
		h.respondWithError(w, http.StatusBadRequest, "Alias required")
		return
	}

	originalURL, err := h.service.ResolveURL(r.Context(), alias)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			h.respondWithError(w, http.StatusNotFound, "URL not found")
			return
		}

		h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func (h *Handler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")
	if alias == "" {
		h.respondWithError(w, http.StatusBadRequest, "Alias required")
		return
	}

	analytics, err := h.service.GetAnalytics(r.Context(), alias)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			h.respondWithError(w, http.StatusNotFound, "URL not found")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.respondWithJSON(w, http.StatusOK, analytics)
}

func (h *Handler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (h *Handler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}
