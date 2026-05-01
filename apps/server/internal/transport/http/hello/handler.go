package hello

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	apphello "github.com/yorukot/netstamp/internal/application/hello"
	"github.com/yorukot/netstamp/internal/transport/http/respond"
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

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetGreeting(r.Context())
	if err != nil {
		respond.Error(w, r, http.StatusServiceUnavailable, "request_cancelled", "request was cancelled")
		return
	}

	respond.JSON(w, http.StatusOK, GreetingResponse{
		Message: result.Message,
		Service: result.Service,
	})
}
