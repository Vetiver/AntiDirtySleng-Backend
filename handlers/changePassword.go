package handlers

import (
	"atnidirtysleng/db"
	"net/http"

	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func getEmailFromToken(tokenString string, user *db.UserChangePassData) (db.UserChangePassData, error) {
	token, err := parseToken(tokenString)
	if err != nil {
		return *user, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email := claims["email"].(string)
		user.Email = email
		return *user, nil
	}

	return *user, fmt.Errorf("Invalid token")
}

func (h *BaseHandler) ChangePassword(c *gin.Context) {
	var userData db.UserChangePassData
	if err := c.BindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := getEmailFromToken(userData.Token, &userData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unexpected error occurred"})
		return
	}

	er := h.db.ChangePassword(user.Email, userData.Password)

	if er != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": er.Error()})
		return
	}
}

