package utils

import (
	"errors"
	"fmt"
	"john/config"
	"john/types"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

func GenerateToken(user types.Employee) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": user.Username,
        "id":       user.ID,
        "isAdmin":  user.IsAdmin,  // Changed to match struct
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })
    return token.SignedString(config.JwtSecret)
}

func ValidateToken(tokenString string) (types.Employee, error) {
	
	var emp types.Employee
	if tokenString == "" {
		return emp, errors.New("authorization token required")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.JwtSecret, nil
	})
	if err != nil {
		return emp, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return emp, errors.New("token expired")
			}
		}

		err := mapstructure.Decode(claims, &emp)
		return emp, err
	}
	return emp, errors.New("invalid token")
}