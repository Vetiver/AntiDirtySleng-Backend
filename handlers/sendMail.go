package handlers

import (
	"atnidirtysleng/db"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func generateRandomCode() int {
	rand.NewSource(time.Now().UnixNano())
	return rand.Intn(90000) + 10000
}

func generateJWT(email string, code int) (string, error) {
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"code":  code,
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *BaseHandler) sendConfirmationEmail(reqData *db.User, code int) string {
	err := godotenv.Load()
	reqData.ConfirmCode = code
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	jwtCode, err := generateJWT(reqData.Email, reqData.ConfirmCode)
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
	m.SetBody("text/html", fmt.Sprintf("Спасибо за регистрацию, вот ваша ссылка на подтверждение: <a href=\"http://localhost:3000/ru/confirmRegister?code=%v\">http://localhost:3000/ru/confirmRegister?code=%v</a>", jwtCode, jwtCode))
	h.Code[reqData.Email] = reqData
	log.Printf("Code for user %s: %d\n", reqData.Email, code)
	d := gomail.NewDialer(smtpName, port, emailAdress, emailPass)

	if err := d.DialAndSend(m); err != nil {
		return ""
	}

	return jwtCode
}

func (h *BaseHandler) SendMail(c *gin.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()
	var reqData *db.User
	code := generateRandomCode()
	if err := c.BindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	jw := h.sendConfirmationEmail(reqData, code)
	c.JSON(http.StatusOK, gin.H{
		"result": jw,
	})
}
