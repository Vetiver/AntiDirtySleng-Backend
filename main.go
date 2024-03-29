package main

import (
	"atnidirtysleng/db"
	"atnidirtysleng/handlers"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	m, err := migrate.New(
		"file://migrations",
		"postgres://postgres:228@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	pool := db.DbStart(databaseURL)

	db := db.NewDB(pool)
	handler := handlers.NewBaseHandler(db)
	r := gin.Default()
	r.GET("/getAllUsers", func(c *gin.Context) {
		handler.GetAllUsers(c)
	})
	// v1 := r.Group("/user")
	// {
	// 	v1.POST("/sendMail", func(c *gin.Context) {
	// 		handler.SendMail(c)
	// 	})
	// 	v1.POST("/regiter", func(c *gin.Context) {
	// 		handler.RegisterUser(c)
	// 	})
	// 	v1.GET("/login", func(c *gin.Context) {
	// 		handler.LoginUser(c)
	// 	})
	// }

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
