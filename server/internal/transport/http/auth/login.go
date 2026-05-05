package auth

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

func (h *Handler) login(_ context.Context, _ *loginInput) (*loginOutput, error) {
	return nil, huma.Error501NotImplemented("login is not implemented")
}

type loginInput struct {
	Body loginInputBody
}

type loginOutput struct {
	Body loginOutputBody
}

type loginInputBody struct {
	Email    string `json:"email" format:"email" required:"true"`
	Password string `json:"password" required:"true"`
}

type loginOutputBody struct {
	TokenType   string `json:"tokenType"`
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
}
