package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"john/utils"
	"log"
	"net/http"
	"strings"
)
// auth_middleware.go
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for health checks, OPTIONS, and login mutation
		if r.URL.Path == "/health" || r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// For GraphQL requests, check if it's a login mutation
		if r.Method == "POST" && r.URL.Path == "/graphql" {
			bodyBytes, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			
			var reqBody struct {
				Query     string `json:"query"`
				Operation string `json:"operationName"`
			}
			if err := json.Unmarshal(bodyBytes, &reqBody); err == nil {
				if strings.Contains(strings.ToLower(reqBody.Query), "mutation") && 
				   strings.Contains(strings.ToLower(reqBody.Query), "login") {
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			respondWithError(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		token := tokenParts[1]
		claims, err := utils.ValidateToken(token)
		if err != nil {
			log.Printf("Token validation failed: %v", err)
			respondWithError(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func respondWithError(w http.ResponseWriter, s string, i int) {
	panic("unimplemented")
}