package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JwtSecret []byte

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	JwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(JwtSecret) == 0 {
		JwtSecret = []byte("GraphQL") // fallback
	}
}
