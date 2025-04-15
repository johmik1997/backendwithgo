package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	db "john/database"
	"john/middleware"
	"john/schema"

	"github.com/rs/cors"
)

func main() {
	// Verify database connection before starting
	if err := db.DB.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
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

	// CORS configuration
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://backendwithgo.onrender.com", "http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposedHeaders:   []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400,
	})

	// Chain middleware
	handler := corsHandler.Handler(
		middleware.LoggingMiddleware(
			middleware.RecoveryMiddleware(mux),
	),
	)
	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	// Configure server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server running on :%s", port)
	log.Fatal(server.ListenAndServe())
}