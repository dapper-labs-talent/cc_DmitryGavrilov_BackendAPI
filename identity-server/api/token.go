package api

import (
	"time"

	"github.com/dapper-labs/identity-server/model"
	"github.com/golang-jwt/jwt/v4"
)

type jwtToken struct {
	jwt.StandardClaims
	Email string
}

func newJwtToken(user *model.User, expiration time.Duration, secret string) (string, error) {
	claims := &jwtToken{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiration).Unix(),
		},
		Email: user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
