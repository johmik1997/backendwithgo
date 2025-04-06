// types/employee.go
package types

import (
	"time"
)

type Employee struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		Password  string    `json:"password"`
		IsAdmin   bool      `json:"isadmin"`
		CreatedAt time.Time `json:"createdat"`
	
}