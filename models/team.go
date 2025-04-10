package models

import (
	//"database/sql"
	"fmt"
	db "john/database"
	
)

type TeamMember struct {
    ID           int       `json:"id"`
    EmpID        int       `json:"empId"`
    EmpName      string    `json:"empName"`
    Department   string    `json:"department"`
    Experience   int       `json:"experience"`
    Address      string    `json:"address"`
    Birthdate    string    `json:"birthdate"`  // Changed to string to match query
    EmployePhoto string    `json:"employePhoto"`
    CreatedAt    string    `json:"createdAt"`  // Changed to string to match query
}
func GetAllTeamMembers() ([]map[string]interface{}, error) {
    query := `
        SELECT 
            ed.id, 
            ed.emp_id AS "empId", 
            ed.emp_name AS "empName", 
            ed.department, 
            ed.experience, 
            ed.address, 
            TO_CHAR(ed.birthdate, 'YYYY-MM-DD') AS "birthdate",
            ed.employephoto AS "employePhoto",
            TO_CHAR(e.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') AS "createdAt"
        FROM employee_details ed
        JOIN employees e ON ed.emp_id = e.id
        ORDER BY e.created_at DESC
    `
    
    rows, err := db.DB.Query(query)
    if err != nil {
        return nil, fmt.Errorf("database query failed: %v", err)
    }
    defer rows.Close()

    var members []map[string]interface{}
    cols, _ := rows.Columns()

    for rows.Next() {
        columns := make([]interface{}, len(cols))
        columnPointers := make([]interface{}, len(cols))
        for i := range columns {
            columnPointers[i] = &columns[i]
        }

        if err := rows.Scan(columnPointers...); err != nil {
            return nil, fmt.Errorf("row scan failed: %v", err)
        }

        m := make(map[string]interface{})
        for i, colName := range cols {
            val := columnPointers[i].(*interface{})
            m[colName] = *val
        }
        members = append(members, m)
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %v", err)
    }
    
    return members, nil
}