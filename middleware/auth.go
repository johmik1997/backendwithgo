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
func respondWithError(w http.ResponseWriter, message string, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "errors": []map[string]interface{}{
            {
                "message": message,
            },
        },
    })
}
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Skip auth for health checks and OPTIONS
        if r.URL.Path == "/health" || r.Method == "OPTIONS" {
            next.ServeHTTP(w, r)
            return
        }

        // For GraphQL requests, parse the query first
        if r.Method == "POST" && r.URL.Path == "/graphql" {
            // Read and restore the request body
            bodyBytes, _ := io.ReadAll(r.Body)
            r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
            
            var reqBody struct {
                Query string `json:"query"`
            }
            
            if err := json.Unmarshal(bodyBytes, &reqBody); err == nil {
                // Allow introspection queries without auth
                if strings.Contains(reqBody.Query, "__schema") || 
                   strings.Contains(reqBody.Query, "__typename") {
                    next.ServeHTTP(w, r)
                    return
                }
                
                // Allow login mutation without auth
                if strings.Contains(strings.ToLower(reqBody.Query), "mutation") && 
                   strings.Contains(strings.ToLower(reqBody.Query), "login") {
                    next.ServeHTTP(w, r)
                    return
                }
            }
        }

        // Require auth for all other requests
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            writeGraphQLError(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            writeGraphQLError(w, "Invalid Authorization header format", http.StatusUnauthorized)
            return
        }

        token := tokenParts[1]
        claims, err := utils.ValidateToken(token)
        if err != nil {
            log.Printf("Token validation failed: %v", err)
            writeGraphQLError(w, "Invalid or expired token", http.StatusUnauthorized)
            return
        }

        // Add claims to context
        ctx := context.WithValue(r.Context(), "user", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func writeGraphQLError(w http.ResponseWriter, message string, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "errors": []map[string]interface{}{
            {
                "message": message,
                "extensions": map[string]interface{}{
                    "code": "UNAUTHENTICATED",
                },
            },
        },
    })
}