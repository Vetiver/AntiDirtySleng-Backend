package handlers

import (
	"atnidirtysleng/db"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func generateJWTConfirmEmail(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *BaseHandler) sendChangeConfirmationEmail(reqData *db.UserEmailData) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	jwtCode, err := generateJWTConfirmEmail(reqData.Email)
	if err != nil {
		log.Fatal("Error jwt generate")
	}
	emailAdress := os.Getenv("EMAIL_ADDRESS")
	emailPass := os.Getenv("EMAIL_PASSWORDCONF")
	smtpName := os.Getenv("SMTP")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	m := gomail.NewMessage()
	m.SetHeader("From", emailAdress)
	m.SetHeader("To", reqData.Email)
	m.SetHeader("Subject", "Confirmation Email")
	m.SetBody("text/html", fmt.Sprintf("Решил поменять пароль братик? Вот твоя ссылочка, дальше ты знаешь что делать: <a href=\"http://localhost:3000/changePassword?code=%d\">http://localhost:8000/changePassword?code=%d</a>", jwtCode, jwtCode))
	log.Printf("Code for user %s: %d\n", reqData.Email)
	d := gomail.NewDialer(smtpName, port, emailAdress, emailPass)

	if err := d.DialAndSend(m); err != nil {
		return ""
	}

	return jwtCode
}

func (h *BaseHandler) SendChangeMail(c *gin.Context) {
	var reqData *db.UserEmailData
	if err := c.BindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	jw := h.sendChangeConfirmationEmail(reqData)
	c.JSON(http.StatusOK, gin.H{
		"result": jw,
	})
}
