package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func createJWTToken(id uuid.UUID) (string, error) {
	JWTToken := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{

		"id": id,
	})

	JWTTokenString
}
