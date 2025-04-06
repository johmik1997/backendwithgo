package models

import (
	"database/sql"
	"errors"
	db "john/database"
	"john/security"
	"john/types"
	"log"
)

func CreateEmployee(username, password string) (*types.Employee, error) {
	hashed, err := security.HashPassword(password)

	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	var emp types.Employee
	err = db.DB.QueryRow(
		"INSERT INTO employees (username, password) VALUES ($1, $2) RETURNING id, username",
		username, hashed,
	).Scan(&emp.ID, &emp.Username)
	
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func GetEmployeeByUsername(username string) (*types.Employee, error) {
	var emp types.Employee
	err := db.DB.QueryRow(
		"SELECT id, username, password, admin, created_at FROM employees WHERE username = $1",
		username,
	).Scan(&emp.ID, &emp.Username, &emp.Password, &emp.IsAdmin, &emp.CreatedAt)
	
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &emp, nil
}

func VerifyEmployeeCredentials(username, password string) (*types.Employee, error) {
    log.Printf("Verifying credentials for: %s", username)
    
    emp, err := GetEmployeeByUsername(username)
    if err != nil {
        log.Printf("User not found: %v", err)
        return nil, err
    }

    log.Printf("Stored hash: %s", emp.Password)
    log.Printf("Input password: %s", password)
    
    if !security.CheckPasswordHash(password, emp.Password) {
        log.Printf("Password mismatch for user %s", username)
        return nil, errors.New("invalid credentials")
    }

    return emp, nil
}