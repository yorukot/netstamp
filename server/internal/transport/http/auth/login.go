package auth

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	appauth "github.com/yorukot/netstamp/internal/application/auth"
)

func (h *Handler) login(ctx context.Context, input *loginInput) (*loginOutput, error) {
	result, err := h.service.Login(ctx, appauth.LoginInput{
		Email:    input.Body.Email,
		Password: input.Body.Password,
	})

	if err != nil {
		switch {
		case errors.Is(err, appauth.ErrCredentialsInvalid):
			return nil, huma.Error401Unauthorized("invalid email or password")
		default:
			return nil, huma.Error500InternalServerError("login failed")
		}
	}

	return &loginOutput{
		Body: loginOutputBody{
			User: userResponse{
				ID:    result.UserID,
				Email: result.Email,
			},
			TokenType:   result.TokenType,
			AccessToken: result.AccessToken,
			ExpiresIn:   result.ExpiresIn,
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
	Email    string `json:"email" format:"email" required:"true"`
	Password string `json:"password" required:"true"`
}

type loginOutputBody struct {
	User        userResponse `json:"user"`
	TokenType   string       `json:"tokenType" example:"Bearer"`
	AccessToken string       `json:"accessToken"`
	ExpiresIn   int          `json:"expiresIn" example:"43200" doc:"Access token lifetime in seconds."`
}
