package main

import (
	"log"
	"net/http"
	"time"
"encoding/json"
	db "john/database"
	"john/middleware"
	"john/schema"

	"github.com/rs/cors"
)

func main() {
	defer db.Close()

	mux := http.NewServeMux()

	// GraphQL endpoint
	mux.Handle("/graphql", middleware.AuthMiddleware(schema.GraphQLHandler()))

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
			log.Printf("Failed to encode health check response: %v", err)
		}
	})

	// main.go
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://backendwithgo.onrender.com", "http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposedHeaders:   []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400,
	})
	// Chain middleware with proper ordering
	handler := corsHandler.Handler(
		middleware.LoggingMiddleware(
			middleware.RecoveryMiddleware(mux),
		),
	)

	// Configure server with timeouts
	server := &http.Server{
		Addr:         ":8082",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Server running on :8082")
	log.Fatal(server.ListenAndServe())
}