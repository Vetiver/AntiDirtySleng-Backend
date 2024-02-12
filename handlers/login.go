package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func createJWTToken(id uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{

		"id": id,
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
