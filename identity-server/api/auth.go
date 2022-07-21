package api

import (
	"context"
	"net/http"
	"strings"
)

const (
	jwtContextKey = "jwtToken"
)

func (api *API) TokenAuth(w http.ResponseWriter, r *http.Request) (context.Context, error) {

	stoken := r.Header.Get("X-Authentication-Token")
	if stoken == "" {
		stoken = extractBearerToken(r.Header.Get("Authorization"))
	}

	if stoken == "" {
		return nil, unauthorizedError("authorization failed, token is missing")
	}

	token, err := api.parseJwtToken(stoken)
	if err != nil {
		return nil, unauthorizedError("authorization failed, reading token failed")
	}

	if !token.Valid {
		return nil, unauthorizedError("authorization failed, token is invalid")
	}
	ctx := r.Context()
	return context.WithValue(ctx, jwtContextKey, token), nil
}

func extractBearerToken(headerValue string) string {

	if headerValue == "" {
		return ""
	}

	// regex is a better solution to match the Bearer token,
	// due to insufficient time I am going to use just a hasPrefix and a split functions
	if !strings.HasPrefix(headerValue, "Bearer ") {
		return ""
	}

	parts := strings.Split(headerValue, " ")
	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}

func bearerAuthFailed(w http.ResponseWriter) {
	w.Header().Add("WWW-Authenticate", "Bearer")
	w.WriteHeader(http.StatusUnauthorized)
}
