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

// func AuthMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Skip auth for health checks and OPTIONS
// 		if r.URL.Path == "/health" || r.Method == "OPTIONS" {
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		// For GraphQL requests, check if it's a login mutation
// 		if r.Method == "POST" && r.URL.Path == "/graphql" {
// 			bodyBytes, _ := io.ReadAll(r.Body)
// 			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			
// 			var reqBody struct {
// 				Query string `json:"query"`
// 			}
// 			if err := json.Unmarshal(bodyBytes, &reqBody); err == nil {
// 				// Skip auth for login mutation
// 				if isLoginMutation(reqBody.Query) {
// 					next.ServeHTTP(w, r)
// 					return
// 				}
// 			}
// 		}

// 		// Require auth for all other requests
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" {
// 			sendAuthError(w, "Authorization header required")
// 			return
// 		}

// 		tokenParts := strings.Split(authHeader, " ")
// 		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
// 			sendAuthError(w, "Invalid Authorization header format")
// 			return
// 		}

// 		token := tokenParts[1]
// 		claims, err := utils.ValidateToken(token)
// 		if err != nil {
// 			log.Printf("Token validation failed: %v", err)
// 			sendAuthError(w, "Invalid or expired token")
// 			return
// 		}

// 		// Add claims to context
// 		ctx := context.WithValue(r.Context(), "user", claims)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

// func isLoginMutation(query string) bool {
// 	query = strings.ToLower(query)
// 	return strings.Contains(query, "mutation") && 
// 	       (strings.Contains(query, "login") || 
// 	        strings.Contains(query, "signin"))
// }

// func sendAuthError(w http.ResponseWriter, message string) {

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusUnauthorized)
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"errors": []map[string]interface{}{
// 			{
// 				"message": message,
// 				"extensions": map[string]string{
// 					"code": "UNAUTHENTICATED",
// 				},
// 			},
// 		},
// 	})
// }

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Skip auth for health checks and OPTIONS
        if r.URL.Path == "/health" || r.Method == "OPTIONS" {
            next.ServeHTTP(w, r)
            return
        }

        // For GraphQL requests, check the operation
        if r.Method == "POST" && r.URL.Path == "/graphql" {
            // Read and restore the request body
            bodyBytes, err := io.ReadAll(r.Body)
            if err != nil {
                writeGraphQLError(w, "Failed to read request body", http.StatusBadRequest)
                return
            }
            r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
            
            var reqBody struct {
                Query         string `json:"query"`
                OperationName string `json:"operationName"`
            }
            
            if err := json.Unmarshal(bodyBytes, &reqBody); err == nil {
                // Skip auth for specific operations
                if shouldSkipAuth(reqBody.Query, reqBody.OperationName) {
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

func shouldSkipAuth(query, operationName string) bool {
    // Convert to lowercase for case-insensitive comparison
    query = strings.ToLower(query)
    operationName = strings.ToLower(operationName)

    // Check for login or register operations
    return strings.Contains(query, "mutation") && 
           (strings.Contains(query, "login") || 
            strings.Contains(query, "register") ||
            operationName == "login" ||
            operationName == "registeruser")
}

func writeGraphQLError(w http.ResponseWriter, message string, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "errors": []map[string]interface{}{
            {
                "message": message,
                "extensions": map[string]string{
                    "code": "UNAUTHENTICATED",
                },
            },
        },
    })
}