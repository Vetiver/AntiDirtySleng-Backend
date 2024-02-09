package handlers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (h BaseHandler) GetAllUsers(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует токен авторизации"})
		return
	}

	token, err := parseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
		return
	}
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
		return
	}
	userID := int(token.Claims.(jwt.MapClaims)["userid"].(int))
	users, err := h.db.GetAllUsers(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Что-то пошло не так"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
