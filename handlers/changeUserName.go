package handlers

import (
	"atnidirtysleng/db"
	"fmt"
	"net/http"

	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

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

func (h BaseHandler) ChangeUsername(c *gin.Context) {
	var changeUsernameRequest db.UserChangeUsernameData
	if err := c.BindJSON(&changeUsernameRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenString := c.GetHeader("Authorization")
	userID, err := getUserIDFromToken(tokenString)
	if err != nil {
		log.Println("Error extracting user ID from token:", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err = h.db.ChangeUsername(userID, changeUsernameRequest.Username)
	if err != nil {
		log.Printf("Error updating username in the database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Username updated successfully"})
}
