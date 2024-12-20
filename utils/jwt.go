package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(email, userID string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Unexpected Signing Method.")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", errors.New("Could not Parse Token.")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return "", errors.New("Invalid Token.")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return "", errors.New("Invalid Token Claims.")
	}

	//email := claims["email"].(string)
	userID := claims["userID"].(string)
	return userID, nil
}
