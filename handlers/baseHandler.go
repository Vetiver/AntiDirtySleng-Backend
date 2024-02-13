package handlers

import (
	"atnidirtysleng/db"
	"github.com/dgrijalva/jwt-go"
	"os"
	"sync"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

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
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
