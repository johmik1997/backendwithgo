package models

import (
	"fmt"
	db "john/database"
	"time"
)
	 type EmpDetails struct {
		ID           int       `json:"id"`
		EmpID        int       `json:"empId"`
		EmpName      string    `json:"empName"`  // Changed from emp_name
		Department   string    `json:"department"`
		Experience   int       `json:"experience"`
		Address      string    `json:"address"`
		Birthdate    time.Time `json:"birthdate"`
		EmployePhoto string    `json:"employePhoto"`
	}
	func AddEmployee(employee *EmpDetails) error {
		// Begin transaction
		tx, err := db.DB.Begin()
		if err != nil {
			return err
		}
	
		// Insert into employee_details
		var id int
		err = tx.QueryRow(`
			INSERT INTO employee_details 
			(emp_id, emp_name, department, experience, address, birthdate, employePhoto)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id
		`,
			employee.EmpID, 
			employee.EmpName, 
			employee.Department, 
			employee.Experience, 
			employee.Address, 
			employee.Birthdate,
			employee.EmployePhoto,
		).Scan(&id)
		
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("could not insert employee details: %v", err)
		}
	
		employee.ID = id
		return tx.Commit()
	}
	func GetAllEmpDetails() ([]EmpDetails, error) {
		rows, err := db.DB.Query(`
			SELECT id, emp_id AS empId, 
				   emp_name AS empName,  <!-- Using SQL aliases -->
				   department, experience, 
				   address, birthdate, employePhoto 
			FROM employee_details
		`)
		// ... rest of the function
		if err != nil {
			return nil, err
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
			)
			if err != nil {
				return nil, err
			}
			employees = append(employees, emp)
		}
		return employees, nil
	}