package hello

import (
	"net/http"

	"github.com/yorukot/netstamp/internal/transport/http/respond"
)

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
