package hello

import (
	"github.com/go-chi/chi/v5"

	apphello "github.com/yorukot/netstamp/internal/application/hello"
)

type Handler struct {
	service *apphello.Service
}

func NewHandler(service *apphello.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/hello", h.Get)
}
