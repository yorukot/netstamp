package auth

import (
	"context"

	"github.com/danielgtaylor/huma/v2"

	httpmiddleware "github.com/yorukot/netstamp/internal/transport/http/middleware"
)

func (h *Handler) me(ctx context.Context, _ *meInput) (*meOutput, error) {
	claims, ok := httpmiddleware.AccessTokenClaimsFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("missing bearer token")
	}

	return &meOutput{
		Body: meOutputBody{
			Authenticated: true,
			User: userResponse{
				ID:    claims.Subject,
				Email: claims.Email,
			},
		},
	}, nil
}

type meInput struct{}

type meOutput struct {
	Body meOutputBody
}

type meOutputBody struct {
	Authenticated bool         `json:"authenticated" example:"true"`
	User          userResponse `json:"user"`
}
