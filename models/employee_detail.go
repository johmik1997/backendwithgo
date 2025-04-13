package models

import (
	"fmt"
	db "john/database"
	"time"
)

type EmpDetails struct {
	ID          int       `json:"id" db:"id"`
	EmpID       int       `json:"empId" db:"emp_id"`
	EmpName     string    `json:"empName" db:"emp_name"`
	Department  string    `json:"department" db:"department"`
	Experience  int       `json:"experience" db:"experience"`
	Address     string    `json:"address" db:"address"`
	Birthdate   time.Time `json:"birthdate" db:"birthdate"`
	EmployePhoto string   `json:"employePhoto" db:"employee_photo"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}
type EmpUpdateDetails struct {
	ID          int       `json:"id" db:"id"`
	EmpID       int       `json:"empId" db:"emp_id"`
	EmpName     string    `json:"empName" db:"emp_name"`
	Department  string    `json:"department" db:"department"`
	Experience  int       `json:"experience" db:"experience"`
	Address     string    `json:"address" db:"address"`
	Birthdate   time.Time `json:"birthdate" db:"birthdate"`
	EmployePhoto string   `json:"employePhoto" db:"employee_photo"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

func AddEmployee(employee *EmpDetails) error {
	// Validate required fields
	if employee.EmpName == "" || employee.Department == "" || employee.Address == "" {
		return fmt.Errorf("empName, department, and address are required")
	}

	// Set default values if needed
	if employee.Birthdate.IsZero() {
		employee.Birthdate = time.Now()
	}
	if employee.CreatedAt.IsZero() {
		employee.CreatedAt = time.Now()
	}

	// Begin transaction
	tx, err := db.DB.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}

	// Insert into employee_details
	query := `
		INSERT INTO employee_details 
		(emp_id, emp_name, department, experience, address, birthdate, employee_photo, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err = tx.QueryRow(
		query,
		employee.EmpID,
		employee.EmpName,
		employee.Department,
		employee.Experience,
		employee.Address,
		employee.Birthdate,
		employee.EmployePhoto,
		employee.CreatedAt,
	).Scan(&employee.ID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not insert employee details: %v", err)
	}

	return tx.Commit()
}

func GetAllEmpDetails() ([]EmpDetails, error) {
	query := `
		SELECT id, emp_id, emp_name, department, 
			   experience, address, birthdate, employee_photo, created_at
		FROM employee_details
		ORDER BY created_at DESC
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not query employee details: %v", err)
	}
	defer rows.Close()

	var employees []EmpDetails
	for rows.Next() {
		var emp EmpDetails
		err := rows.Scan(
			&emp.ID,
			&emp.EmpID,
			&emp.EmpName,
			&emp.Department,
			&emp.Experience,
			&emp.Address,
			&emp.Birthdate,
			&emp.EmployePhoto,
			&emp.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan employee details: %v", err)
		}
		employees = append(employees, emp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return employees, nil
}


func UpdateEmployee(emp *EmpUpdateDetails) error {
	query := `
		UPDATE employee_details
		SET emp_name = $1,
			department = $2,
			experience = $3,
			address = $4,
			birthdate = $5,
			employee_photo = $6
		WHERE id = $7
	`

	_, err := db.DB.Exec(
		query,
		emp.EmpName,
		emp.Department,
		emp.Experience,
		emp.Address,
		emp.Birthdate,
		emp.EmployePhoto,
		emp.ID,
	)

	if err != nil {
		return fmt.Errorf("could not update employee: %v", err)
	}

	return nil
}
func DeleteEmployee(id int) error {
    result, err := db.DB.Exec("DELETE FROM employee_details WHERE id = $1", id)
    if err != nil {
        return fmt.Errorf("database error: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("could not verify deletion: %v", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("no employee found with id %d", id)
    }
    
    return nil
}

func CheckDBConnection() error {
    return db.DB.Ping()
}