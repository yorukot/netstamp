package team

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"

	appteam "github.com/yorukot/netstamp/internal/application/team"
	httpmiddleware "github.com/yorukot/netstamp/internal/transport/http/middleware"
)

func currentUserID(ctx context.Context) (string, error) {
	claims, ok := httpmiddleware.AccessTokenClaimsFromContext(ctx)
	if !ok || claims.Subject == "" {
		return "", huma.Error401Unauthorized("missing bearer token")
	}

	return claims.Subject, nil
}

func mapTeamError(err error, fallback string) error {
	switch {
	case errors.Is(err, appteam.ErrTeamNotFound), errors.Is(err, appteam.ErrMemberNotFound), errors.Is(err, appteam.ErrUserNotFound):
		return huma.Error404NotFound("not found")
	case errors.Is(err, appteam.ErrForbidden):
		return huma.Error403Forbidden("forbidden")
	case errors.Is(err, appteam.ErrTeamSlugAlreadyExists):
		return huma.Error409Conflict("team slug already exists")
	case errors.Is(err, appteam.ErrMemberAlreadyExists):
		return huma.Error409Conflict("team member already exists")
	case errors.Is(err, appteam.ErrLastOwner):
		return huma.Error409Conflict("team must keep an owner")
	case errors.Is(err, appteam.ErrInvalidInput):
		return huma.Error422UnprocessableEntity("invalid team input")
	case errors.Is(err, appteam.ErrInvalidRole):
		return huma.Error422UnprocessableEntity("invalid team member role")
	default:
		return huma.Error500InternalServerError(fallback)
	}
}
