package handlers

import (
	"atnidirtysleng/db"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"fmt"
)

func getEmailAndCodeFromToken(tokenString string, user *db.User) (db.User, error) {
	token, err := parseToken(tokenString)
	if err != nil {
		return *user, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email := claims["email"].(string)
		code := int(claims["code"].(float64))
		user.Email = email
		user.ConfirmCode = code
		return *user, nil
	}

	return *user, fmt.Errorf("Invalid token")
}

func (h BaseHandler) RegisterUser(c *gin.Context) {
	var token *db.Token
	user := &db.User{}
	if err := c.BindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	decodedUser, err := getEmailAndCodeFromToken(token.TokenString, user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}


	storedUser, exists := h.Code[decodedUser.Email]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if storedUser.ConfirmCode != decodedUser.ConfirmCode {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ConfirmCode"})
		return
	}
	decodedUser.Username = storedUser.Username
	decodedUser.Password = storedUser.Password
	fmt.Println(storedUser)
	registeredUser, err := h.db.RegisterUser(decodedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	delete(h.Code, decodedUser.Email)
	c.JSON(http.StatusOK, gin.H{"result": registeredUser})
}
