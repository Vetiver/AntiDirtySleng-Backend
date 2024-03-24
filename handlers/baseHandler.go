package handlers

import (
	"atnidirtysleng/db"
	"fmt"
	"os"
	"sync"

	"github.com/dgrijalva/jwt-go"
)

type UserGet struct {
	Parce []db.User `json:"parce"`
}

type BaseHandler struct {
	db   *db.DB
	Code map[string]*db.User
	mu   sync.Mutex
}

func NewBaseHandler(pool *db.DB) *BaseHandler {
	return &BaseHandler{
		db:   pool,
		Code: make(map[string]*db.User),
	}
}

func parseToken(tokenString string) (*jwt.Token, error) {
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func getUserIDFromToken(tokenString string) (string, error) {
	token, err := parseToken(tokenString)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := claims["id"].(string)
		return id, nil
	}

	return "", fmt.Errorf("Invalid token")
}
