package utils

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"john/config"
	"john/models"
)

func GenerateToken(user models.Employee) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
	})
	return token.SignedString(config.JwtSecret)
}

func ValidateToken(tokenString string) (models.Employee, error) {
	var emp models.Employee
	if tokenString == "" {
		return emp, errors.New("authorization token required")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid method")
		}
		return config.JwtSecret, nil
	})
	if err != nil {
		return emp, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		err := mapstructure.Decode(claims, &emp)
		return emp, err
	}
	return emp, errors.New("invalid token")
}
