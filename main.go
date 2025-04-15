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
	// Verify database connection
	if err := db.DB.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.Handle("/graphql", schema.GraphQLHandler())// Remove AuthMiddleware from here

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// CORS configuration
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://backendwithgo.onrender.com", "https://frontendvue-with-gobackend.onrender.com"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	})

	// Chain middleware with AuthMiddleware first
	handler := corsHandler.Handler(
		middleware.AuthMiddleware(
			middleware.LoggingMiddleware(
				middleware.RecoveryMiddleware(mux),
			),
		),
	)

	// Server configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

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