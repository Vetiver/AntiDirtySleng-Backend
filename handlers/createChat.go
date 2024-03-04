package handlers

import (
	"atnidirtysleng/db"
	"net/http"

	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func getOwnerFromToken(tokenString string) (string, error) {
	token, err := parseToken(tokenString)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		owner := claims["id"].(string)
		return owner, nil
	}

	return "", fmt.Errorf("Invalid token")
}

func (h BaseHandler) CreateChat(c *gin.Context) {
	var chat db.Chat

	log.Println("Start processing CreateChat handler")

	log.Println("Received request body:", c.Request.Body)

	if err := c.ShouldBindJSON(&chat); err != nil {
		log.Println("Error binding JSON:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenString := c.GetHeader("Authorization")
	ownerStr, err := getOwnerFromToken(tokenString)
	if err != nil {
		log.Println("Error extracting owner from token:", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	owner, err := uuid.Parse(ownerStr)
	if err != nil {
		log.Println("Error parsing owner UUID:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse owner UUID"})
		return
	}

	chat.Owner = owner

	err = h.db.CreateChat(&chat)
	if err != nil {
		log.Println("Error creating chat:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	err = h.db.AddUserToChat(ownerStr, chat.ChatId)
	if err != nil {
		log.Println("Error adding user to chat:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to chat"})
		return
	}

	log.Println("CreateChat handler processed successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Chat created successfully"})
}
