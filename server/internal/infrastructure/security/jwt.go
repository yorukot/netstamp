package security

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"

	appauth "github.com/yorukot/netstamp/internal/application/auth"
)

type JWTIssuer struct {
	secret []byte
	ttl    time.Duration
	now    func() time.Time
}

type accessTokenClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewJWTIssuer(secret string, ttl time.Duration) *JWTIssuer {
	return &JWTIssuer{
		secret: []byte(secret),
		ttl:    ttl,
		now:    time.Now,
	}
}

func (i *JWTIssuer) IssueAccessToken(ctx context.Context, input appauth.AccessTokenInput) (appauth.IssuedToken, error) {
	if err := ctx.Err(); err != nil {
		return appauth.IssuedToken{}, err
	}

	now := i.now().UTC()
	expiresAt := now.Add(i.ttl)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims{
		Email: input.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   input.Subject,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	})

	value, err := token.SignedString(i.secret)
	if err != nil {
		return appauth.IssuedToken{}, err
	}

	return appauth.IssuedToken{
		Value:     value,
		TokenType: "Bearer",
		ExpiresIn: int(i.ttl.Seconds()),
	}, nil
}
