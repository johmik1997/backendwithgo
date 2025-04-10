package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"john/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for health checks and OPTIONS
		if r.URL.Path == "/health" || r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
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
			respondWithError(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		log.Printf("Authenticated user: %s (ID: %d)", claims.Username, claims.ID)
		
		// Add claims to context
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}