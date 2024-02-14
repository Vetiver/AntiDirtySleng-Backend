package handlers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	userid, err := uuid.Parse(token.Claims.(jwt.MapClaims)["userid"].(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ошибка при распарсе токена"})
		return
	}
	users, err := h.db.GetAllUsers(userid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Что-то пошло не так"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
