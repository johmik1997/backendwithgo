package main

import (
	"log"
	"net/http"

	//"john/config"
	db "john/database"
	"john/handlers"
	"john/middleware"
	"john/schema"

	"github.com/rs/cors"
)

func main() {
	defer db.Close()

	mux := http.NewServeMux()

	mux.Handle("/graphql", middleware.AuthMiddleware(schema.GraphQLHandler()))
	mux.HandleFunc("/login", handlers.LoginHandler)
	//http.HandleFunc("/register", handlers.RegisterHandler)
	// Set CORS headers
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust as needed
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// Call the RegisterHandler
		handlers.RegisterHandler(w, r)
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(middleware.LoggingMiddleware(mux))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
