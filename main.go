package main

import (
	"log"
	"net/http"

	db "john/database"
	"john/middleware"
	"john/schema"

	"github.com/rs/cors"
)

func main() {
	defer db.Close()

	mux := http.NewServeMux()

	// GraphQL endpoint with auth middleware for protected operations
	mux.Handle("/graphql", middleware.AuthMiddleware(schema.GraphQLHandler()))
	

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// CORS configuration
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(middleware.LoggingMiddleware(mux))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}