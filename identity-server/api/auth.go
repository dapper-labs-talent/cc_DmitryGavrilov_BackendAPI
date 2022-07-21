package api

import (
	"net/http"
	"strings"
)

func (api *API) BearerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// we should expect to get a token value passed in the x-authentication-token header (according to test task description)
		stoken := r.Header.Get("X-Authentication-Token")
		if stoken == "" {
			stoken = extractBearerToken(r.Header.Get("Authorization"))
		}

		if stoken == "" {
			bearerAuthFailed(w)
			return
		}

		token, err := api.parseJwtToken(stoken)
		if err != nil || !token.Valid {
			bearerAuthFailed(w)
			return
		}

		next.ServeHTTP(w, r)
	})
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
