package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func init() {
	connStr := os.Getenv("DATABASE_URL") // Use Render's full connection string

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Connected to PostgreSQL")

	seedDatabase()
}

func Close() {
	DB.Close()
}
