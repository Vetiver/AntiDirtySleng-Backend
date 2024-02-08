package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

type DB struct {
	pool *pgxpool.Pool
}

type User struct {
	UserId      uuid.UUID `json:"id"`
	Username    string    `json:"name"     binding:"required"`
	IsAdmin     bool      `json:"isAdmin"`
	Email       string    `json:"email"    binding:"required"`
	Password    string    `json:"password" binding:"required,min=8"`
	Description string    `json:"descriprion"`
	Avatar      string    `json:"avatar"`
	ConfirmCode int
}

type UserLoginData struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Token struct {
	TokenString string `json:"accessToken"`
}

func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{
		pool: pool,
	}
}


func DbStart(baseUrl string) *pgxpool.Pool {
	urlExample := baseUrl
	dbpool, err := pgxpool.New(context.Background(), string(urlExample))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v", err)
		os.Exit(1)
	}
	return dbpool
}

func (db DB) userExists(userID uuid.UUID) (bool, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return false, fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	var exists bool
	err = conn.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)", userID).
		Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %v", err)
	}

	return exists, nil
}

func (db DB) GetAllUsers(userID uuid.UUID) ([]User, error) {
    exists, err := db.userExists(userID)
    if err != nil {
        return nil, err
    }

    if exists == false {
        return nil, fmt.Errorf("user with ID %s does not exist", userID.String())
    }

    conn, err := db.pool.Acquire(context.Background())
    if err != nil {
        return nil, fmt.Errorf("unable to acquire a database connection: %v", err)
    }
    defer conn.Release()

    rows, err := conn.Query(context.Background(),
        `SELECT id, name, email, description, avatar FROM users`)
    if err != nil {
        return nil, fmt.Errorf("unable to retrieve data from database: %v", err)
    }
    defer rows.Close()

    var data []User
    for rows.Next() {
        var d User
        err = rows.Scan(&d.UserId,&d.Username, &d.Email, &d.Description, &d.Avatar )
        if err != nil {
            return nil, fmt.Errorf("unable to scan row: %v", err)
        }
        data = append(data, d)
    }
    return data, err
}