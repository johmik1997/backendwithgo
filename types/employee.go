// types/employee.go
package types

import (
	"time"
)

type Employee struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Password  string    `json:"password"`
    IsAdmin   bool      `json:"isAdmin"`  // Changed to match database
    CreatedAt time.Time `json:"createdAt"`
}