// package models

// import (
// 	//"database/sql"
// 	"fmt"
// 	db "john/database"
	
// )

// type TeamMember struct {
//     ID           int       `json:"id"`
//     EmpID        int       `json:"empId"`
//     EmpName      string    `json:"empName"`
//     Department   string    `json:"department"`
//     Experience   int       `json:"experience"`
//     Address      string    `json:"address"`
//     Birthdate    string    `json:"birthdate"`  // Changed to string to match query
//     EmployePhoto string    `json:"employePhoto"`
//     CreatedAt    string    `json:"createdAt"`  // Changed to string to match query
// }
// func GetAllTeamMembers() ([]TeamMember, error) {
//     query := `
//         SELECT 
//             ed.id, 
//             ed.emp_id AS "empId", 
//             ed.emp_name AS "empName", 
//             ed.department, 
//             ed.experience, 
//             ed.address, 
//             TO_CHAR(ed.birthdate, 'YYYY-MM-DD') AS "birthdate",
//             ed.employe_photo AS "employePhoto",
//             TO_CHAR(e.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') AS "createdAt"
//         FROM employee_details ed
//         JOIN employees e ON ed.emp_id = e.id
//         ORDER BY e.created_at DESC
//     `
    
//     rows, err := db.DB.Query(query)
//     if err != nil {
//         return nil, fmt.Errorf("database query failed: %v", err)
//     }
//     defer rows.Close()

//     var members []TeamMember
//     for rows.Next() {
//         var m TeamMember
//         err := rows.Scan(
//             &m.ID,
//             &m.EmpID,
//             &m.EmpName,
//             &m.Department,
//             &m.Experience,
//             &m.Address,
//             &m.Birthdate,
//             &m.EmployePhoto,
//             &m.CreatedAt,
//         )
//         if err != nil {
//             return nil, fmt.Errorf("row scan failed: %v", err)
//         }
//         members = append(members, m)
//     }
    
//     if err := rows.Err(); err != nil {
//         return nil, fmt.Errorf("rows error: %v", err)
//     }
    
//     return members, nil
// }