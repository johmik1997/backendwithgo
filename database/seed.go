package db

import (
	"john/security"
	"log"
)

func seedDatabase() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS employees (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(100) NOT NULL,
			admin BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatal("Error creating employees table:", err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS employee_details (
			emp_id INTEGER PRIMARY KEY REFERENCES employees(id),
			emp_title VARCHAR(100),
			address TEXT
		)
	`)
	if err != nil {
		log.Fatal("Error creating employee_details table:", err)
	}

	// Insert employees with hashed passwords
	users := []struct {
		username string
		password string
		admin    bool
	}{
		{"Gordon", "64382", true},
		{"Nick", "7845", false},
		{"John", "1234", false},
	}

	for _, u := range users {
		hashed, err := security.HashPassword(u.password)
		if err != nil {
			log.Fatalf("Error hashing password for %s: %v", u.username, err)
		}
		_, err = DB.Exec(`
			INSERT INTO employees (username, password, admin)
			VALUES ($1, $2, $3)
			ON CONFLICT (username) DO NOTHING
		`, u.username, hashed, u.admin)
		if err != nil {
			log.Println("Insert employee failed:", err)
		}
	}

	// Optional: seed one employee detail row
	_, err = DB.Exec(`
		INSERT INTO employee_details (emp_id, emp_title, address)
		VALUES (1, 'Software Engineer', 'San Jose')
		ON CONFLICT (emp_id) DO NOTHING
	`)
	if err != nil {
		log.Println("Insert employee_details failed:", err)
	}
}
