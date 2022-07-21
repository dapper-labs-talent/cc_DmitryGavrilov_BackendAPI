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

func (api *API) parseJwtToken(token string) (*jwt.Token, error) {

	p := jwt.Parser{ValidMethods: []string{jwt.SigningMethodHS256.Name}}
	ptoken, err := p.ParseWithClaims(token, &jwtToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(api.config.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	return ptoken, nil
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
