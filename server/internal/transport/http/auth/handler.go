package auth

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	appauth "github.com/yorukot/netstamp/internal/application/auth"
)

type Handler struct {
	service *appauth.Service
}

func NewHandler(service *appauth.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "registerUser",
		Method:        http.MethodPost,
		Path:          "/auth/register",
		DefaultStatus: http.StatusCreated,
		Summary:       "Register user",
		Tags:          []string{"Auth"},
		Errors:        []int{http.StatusUnprocessableEntity, http.StatusConflict, http.StatusInternalServerError},
	}, h.register)

	huma.Register(api, huma.Operation{
		OperationID: "loginUser",
		Method:      http.MethodPost,
		Path:        "/auth/login",
		Summary:     "Login user",
		Tags:        []string{"Auth"},
		Errors:      []int{http.StatusNotImplemented},
	}, h.login)
}
