package auth

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

func (h *Handler) register(ctx context.Context, _ *registerInput) (*registerOutput, error) {
	result, err := h.service.GetGreeting(ctx)
	if err != nil {
		return nil, huma.Error503ServiceUnavailable("request was cancelled")
	}

	return &registerOutput{
		Body: registerOutputBody{
			Message: result.Message,
		},
	}, nil
}

type registerInput struct {
	Body registerInputBody
}

type registerOutput struct {
	Body registerOutputBody
}

type registerInputBody struct {
	Username string
	Password string
}

type registerOutputBody struct {
	Message string
	TokenType  string
	AccessToken string
	ExpiresIn int
}
