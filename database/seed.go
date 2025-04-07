package db

import (
	"john/security"
	"log"
)

func seedDatabase() {
	// Create employees table
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS employees (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(100) NOT NULL,
			is_admin BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		log.Fatal("Error creating employees table:", err)
	}

	// Create employee_details table with proper foreign key
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS employee_details (
			id SERIAL PRIMARY KEY,
			emp_id INT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
			emp_name VARCHAR(100) NOT NULL,
			department VARCHAR(100),
			experience INT,
			address TEXT,
			birthdate DATE,
			employePhoto TEXT
		);
	`)
	if err != nil {
		log.Fatal("Error creating employee_details table:", err)
	}

	// Insert employees first (with hashed passwords)
	users := []struct {
		username string
		password string
		isAdmin  bool
	}{
		{"melkam", "64382", true},
		{"mom", "7845", false},
		{"bar", "1234", false},
	}

	var employeeIDs []int
	for _, u := range users {
		// Hash the password
		hashed, err := security.HashPassword(u.password)
		if err != nil {
			log.Fatalf("Error hashing password for %s: %v", u.username, err)
		}

		// Insert employee and get the generated ID
		var id int
		err = DB.QueryRow(`
			INSERT INTO employees (username, password, is_admin)
			VALUES ($1, $2, $3)
			RETURNING id
		`, u.username, hashed, u.isAdmin).Scan(&id)
		if err != nil {
			log.Println("Insert employee failed:", err)
			continue
		}
		employeeIDs = append(employeeIDs, id)
	}

	// Only insert employee details if we have matching employee IDs
	if len(employeeIDs) >= 3 {
		employees := []struct {
			emp_id       int
			emp_name     string
			department   string
			experience   int
			address      string
			birthdate    string
			employePhoto string
		}{
			{employeeIDs[0], "Melkam", "Software Engineer", 5, "San Jose", "1990-04-15", "melkam_photo.jpg"},
			{employeeIDs[1], "Mom", "Product Manager", 3, "Mountain View", "1992-06-22", "mom_photo.jpg"},
			{employeeIDs[2], "Bar", "UX Designer", 4, "Palo Alto", "1991-08-09", "bar_photo.jpg"},
		}

		for _, e := range employees {
			_, err = DB.Exec(`
				INSERT INTO employee_details 
				(emp_id, emp_name, department, experience, address, birthdate, employePhoto)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`, e.emp_id, e.emp_name, e.department, e.experience, e.address, e.birthdate, e.employePhoto)
			if err != nil {
				log.Println("Insert employee details failed:", err)
			}
		}
	}
}