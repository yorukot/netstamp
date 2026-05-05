package auth

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
        OperationID: "registerUser",
        Method:      http.MethodPost,
        Path:        "/auth/register",
        Summary:     "Register user",
        Tags:        []string{"Auth"},
    }, h.register)

    huma.Register(api, huma.Operation{
        OperationID: "loginUser",
        Method:      http.MethodPost,
        Path:        "/auth/login",
        Summary:     "Login user",
        Tags:        []string{"Auth"},
    }, h.login)
}