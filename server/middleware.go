package server

import (
	"net/http"
	"strings"

	"github.com/worming004/TwelveFactorApp/auth"
)

func getAuthMiddleware(jwtWrap auth.JwtWrapper) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			splittedHeader := strings.Split(authHeader, "Bearer ")
			var tokenValue string
			if len(splittedHeader) == 2 {
				tokenValue = strings.TrimSpace(splittedHeader[1])
			} else {
				http.Error(w, "bad format for Authorization header", http.StatusBadRequest)
				return
			}

			_, err := jwtWrap.ValidateToken(tokenValue)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
