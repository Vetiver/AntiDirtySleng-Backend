package handlers

import (
	"atnidirtysleng/db"
	"fmt"
	"net/http"

	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func getUserIDFromTokenForDescription(tokenString string) (string, error) {
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

func (h BaseHandler) ChangeDescription(c *gin.Context) {
	var changeDescriptionRequest db.UserChangeDescriptionData
	if err := c.BindJSON(&changeDescriptionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if changeDescriptionRequest.Description != "" && len(changeDescriptionRequest.Description) > 255 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Description must be no longer than 255 characters"})
		return
	}

	tokenString := c.GetHeader("Authorization")
	userID, err := getUserIDFromTokenForDescription(tokenString)
	if err != nil {
		log.Println("Error extracting user ID from token:", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err = h.db.ChangeDescription(userID, changeDescriptionRequest.Description)
	if err != nil {
		log.Printf("Error updating description in the database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Description updated successfully"})
}
