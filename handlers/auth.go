// package handlers

// import (
// 	"encoding/json"
// 	"net/http"
// 	"john/models"
// 	"john/utils"
// )

// func LoginHandler(w http.ResponseWriter, r *http.Request) {
// 	var user models.Employee
// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	dbUser, err := models.GetEmployeeByUsername(user.Username)
// 	if err != nil || dbUser.Password != user.Password {
// 		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// 		return
// 	}

// 	token, err := utils.GenerateToken(dbUser)
// 	if err != nil {
// 		http.Error(w, "Token generation error", http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(map[string]string{"token": token})
// }

// func RegisterHandler(w http.ResponseWriter, r *http.Request) {
// 	var user models.Employee
// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}
// 	if _, err := models.GetEmployeeByUsername(user.Username); err == nil {
// 		http.Error(w, "User already exists", http.StatusConflict)
// 		return
// 	}

// 	// Insert into database
// 	if err := models.CreateEmployee(user); err != nil {
// 		http.Error(w, "Failed to register", http.StatusInternalServerError)
// 		return
// 	}

// 	// Generate JWT token
// 	tokenString, err := utils.GenerateToken(user)
// 	if err != nil {
// 		http.Error(w, "Failed to create token", http.StatusInternalServerError)
// 		return
// 	}

// 	// Send token in response
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
// }