package main

import (
	"atnidirtysleng/db"
	"atnidirtysleng/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	m, err := migrate.New(
		"file://migrations",
		"postgres://postgres:2778@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	pool := db.DbStart(databaseURL)

	db := db.NewDB(pool)
	handler := handlers.NewBaseHandler(db)
	r := gin.Default()

	// Middleware to ignore CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	r.GET("/getUserInfo", func(c *gin.Context) {
		handler.GetUserInfo(c)
	})
	v1 := r.Group("/auth")
	{
		v1.POST("/sendMail", func(c *gin.Context) {
			handler.SendMail(c)
		})
		v1.POST("/register", func(c *gin.Context) {
			handler.RegisterUser(c)
		})
		v1.POST("/login", func(c *gin.Context) {
			handler.LoginUser(c)
		})
		v1.POST("/refreshToken", func(c *gin.Context) {
			handler.RefreshToken(c)
		})
	}

	r.POST("/confirmEmail", func(c *gin.Context) {
		handler.SendChangeMail(c)
	})

	r.POST("/changePassword", func(c *gin.Context) {
		handler.ChangePassword(c)
	})

	r.POST("/createChat", func(c *gin.Context) {
		handler.CreateChat(c)
	})

	r.POST("/changeDescription", func(c *gin.Context) {
		handler.ChangeDescription(c)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
