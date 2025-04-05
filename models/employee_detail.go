package models

import db "john/database"

type EmpDetails struct {
	EmpId    string `json:"id"`
	EmpTitle string `json:"title"`
	Address  string `json:"address"`
}

func GetAllEmpDetails() ([]EmpDetails, error) {
	rows, err := db.DB.Query("SELECT emp_id, emp_title, address FROM employee_details")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []EmpDetails
	for rows.Next() {
		var d EmpDetails
		if err := rows.Scan(&d.EmpId, &d.EmpTitle, &d.Address); err != nil {
			return nil, err
		}
		details = append(details, d)
	}
	return details, nil
}