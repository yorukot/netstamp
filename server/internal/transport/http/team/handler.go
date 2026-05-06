package team

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	appauth "github.com/yorukot/netstamp/internal/application/auth"
	appteam "github.com/yorukot/netstamp/internal/application/team"
	httpmiddleware "github.com/yorukot/netstamp/internal/transport/http/middleware"
)

type Handler struct {
	service  *appteam.Service
	verifier appauth.TokenVerifier
}

func NewHandler(service *appteam.Service, verifier appauth.TokenVerifier) *Handler {
	return &Handler{
		service:  service,
		verifier: verifier,
	}
}

func (h *Handler) RegisterRoutes(api huma.API) {
	if h.service == nil || h.verifier == nil {
		return
	}

	authMiddleware := httpmiddleware.RequireAuth(h.verifier)
	security := []map[string][]string{{"bearerAuth": {}}}
	middlewares := huma.Middlewares{authMiddleware}

	huma.Register(api, huma.Operation{
		OperationID:   "createTeam",
		Method:        http.MethodPost,
		Path:          "/teams",
		DefaultStatus: http.StatusCreated,
		Summary:       "Create team",
		Tags:          []string{"Teams"},
		Security:      security,
		Middlewares:   middlewares,
		Errors:        []int{http.StatusUnauthorized, http.StatusConflict, http.StatusUnprocessableEntity, http.StatusInternalServerError},
	}, h.createTeam)

	huma.Register(api, huma.Operation{
		OperationID: "listTeams",
		Method:      http.MethodGet,
		Path:        "/teams",
		Summary:     "List teams",
		Tags:        []string{"Teams"},
		Security:    security,
		Middlewares: middlewares,
		Errors:      []int{http.StatusUnauthorized, http.StatusInternalServerError},
	}, h.listTeams)

	huma.Register(api, huma.Operation{
		OperationID: "getTeam",
		Method:      http.MethodGet,
		Path:        "/teams/{ref}",
		Summary:     "Get team",
		Tags:        []string{"Teams"},
		Security:    security,
		Middlewares: middlewares,
		Errors:      []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError},
	}, h.getTeam)

	huma.Register(api, huma.Operation{
		OperationID: "updateTeam",
		Method:      http.MethodPatch,
		Path:        "/teams/{ref}",
		Summary:     "Update team",
		Tags:        []string{"Teams"},
		Security:    security,
		Middlewares: middlewares,
		Errors:      []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound, http.StatusConflict, http.StatusUnprocessableEntity, http.StatusInternalServerError},
	}, h.updateTeam)

	huma.Register(api, huma.Operation{
		OperationID:   "deleteTeam",
		Method:        http.MethodDelete,
		Path:          "/teams/{ref}",
		DefaultStatus: http.StatusNoContent,
		Summary:       "Delete team",
		Tags:          []string{"Teams"},
		Security:      security,
		Middlewares:   middlewares,
		Errors:        []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound, http.StatusInternalServerError},
	}, h.deleteTeam)

	huma.Register(api, huma.Operation{
		OperationID: "listTeamMembers",
		Method:      http.MethodGet,
		Path:        "/teams/{ref}/members",
		Summary:     "List team members",
		Tags:        []string{"Team Members"},
		Security:    security,
		Middlewares: middlewares,
		Errors:      []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError},
	}, h.listMembers)

	huma.Register(api, huma.Operation{
		OperationID:   "addTeamMember",
		Method:        http.MethodPost,
		Path:          "/teams/{ref}/members",
		DefaultStatus: http.StatusCreated,
		Summary:       "Add team member",
		Tags:          []string{"Team Members"},
		Security:      security,
		Middlewares:   middlewares,
		Errors:        []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound, http.StatusConflict, http.StatusUnprocessableEntity, http.StatusInternalServerError},
	}, h.addMember)

	huma.Register(api, huma.Operation{
		OperationID: "updateTeamMemberRole",
		Method:      http.MethodPatch,
		Path:        "/teams/{ref}/members/{user_id}",
		Summary:     "Update team member role",
		Tags:        []string{"Team Members"},
		Security:    security,
		Middlewares: middlewares,
		Errors:      []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound, http.StatusConflict, http.StatusUnprocessableEntity, http.StatusInternalServerError},
	}, h.updateMemberRole)
}
