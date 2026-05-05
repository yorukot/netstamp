package hello

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	apphello "github.com/yorukot/netstamp/internal/application/hello"
)

type Handler struct {
	service *apphello.Service
}

func NewHandler(service *apphello.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "getGreeting",
		Method:      http.MethodGet,
		Path:        "/hello",
		Summary:     "Get greeting",
		Tags:        []string{"Hello"},
		Errors:      []int{http.StatusServiceUnavailable},
	}, h.Get)
}
