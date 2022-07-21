package api

import (
	"context"
	"errors"
	"time"

	"github.com/dapper-labs/identity-server/model"
	"github.com/golang-jwt/jwt/v4"
)

type jwtTokenClaims struct {
	jwt.StandardClaims
	Email string
}

func (api *API) parseJwtToken(token string) (*jwt.Token, error) {

	p := jwt.Parser{ValidMethods: []string{jwt.SigningMethodHS256.Name}}
	ptoken, err := p.ParseWithClaims(token, &jwtTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(api.config.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	return ptoken, nil
}

func newJwtToken(user *model.User, expiration time.Duration, secret string) (string, error) {
	claims := &jwtTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiration).Unix(),
		},
		Email: user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func getClaims(context context.Context) (*jwtTokenClaims, error) {
	token := getJwtToken(context)
	if token == nil {
		return nil, errors.New("cannot find token")
	}

	claims, ok := token.Claims.(*jwtTokenClaims)
	if !ok {
		return nil, errors.New("cannot read jwt token claims")
	}

	return claims, nil
}

func getJwtToken(ctx context.Context) *jwt.Token {
	token := ctx.Value(jwtContextKey)
	if token == nil {
		return nil
	}

	return token.(*jwt.Token)
}
