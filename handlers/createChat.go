package handlers

import (
	"atnidirtysleng/db"
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
	ownerStr, err := getUserIDFromToken(tokenString)
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

	if chat.Users == nil {
		chat.Users = []string{}
	}

	ownerExists := false
	for _, userID := range chat.Users {
		if userID == ownerStr {
			ownerExists = true
			break
		}
	}
	if !ownerExists {
		chat.Users = append(chat.Users, ownerStr)
	}

	userMap := make(map[string]bool)
	for _, userID := range chat.Users {
		if userMap[userID] {
			log.Println("Duplicate user ID:", userID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate user ID"})
			return
		}
		userMap[userID] = true
	}

	err = h.db.CreateChat(&chat)
	if err != nil {
		log.Println("Error creating chat:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	for _, userID := range chat.Users {
		err := h.db.AddUserToChat(userID, chat.ChatId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add users"})
			panic(err)
		}
	}

	log.Println("CreateChat handler processed successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Chat created successfully"})
}
