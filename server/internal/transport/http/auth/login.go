package auth

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

func (h *Handler) login(ctx context.Context, _ *loginInput) (*loginOutput, error) {
	result, err := h.service.GetGreeting(ctx)
	if err != nil {
		return nil, huma.Error503ServiceUnavailable("request was cancelled")
	}

	return &loginOutput{
		Body: loginOutputBody{
			Message: result.Message,
		},
	}, nil
}

type loginInput struct {
	Body loginInputBody
}

type loginOutput struct {
	Body loginOutputBody
}

type loginInputBody struct {
	Username string
	Password string
}

type loginOutputBody struct {
	Message string
	TokenType  string
	AccessToken string
	ExpiresIn int
}

