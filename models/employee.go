package models

import (
	db "john/database"
)

type Employee struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetEmployeeByUsername(username string) (Employee, error) {
	var e Employee
	err := db.DB.QueryRow(`SELECT id, username, password FROM employees WHERE username=$1`, username).
		Scan(&e.ID, &e.Username, &e.Password)
	return e, err
}
func CreateEmployee(user Employee) error {
    _, err := db.DB.Exec("INSERT INTO employees (username, password) VALUES ($1, $2)", user.Username, user.Password)
    return err
}
