package middleware

import (
	"context"
	"net/http"
	"strings"
	"john/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")

		if token != "" {
			if _, err := utils.ValidateToken(token); err == nil {
				ctx := context.WithValue(r.Context(), "token", token)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}
