package db

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	pool *pgxpool.Pool
}

type User struct {
	UserId       uuid.UUID `json:"id"`
	Username     string    `json:"name"     binding:"required"`
	IsAdmin      bool      `json:"isAdmin"`
	Email        string    `json:"email"    binding:"required"`
	Password     string    `json:"password" binding:"required,min=8"`
	Description  string    `json:"description"`
	Avatar       *string   `json:"avatar"`
	ConfirmCode  int       `json:"confirmCode"`
	RefreshToken string    `json:"refreshToken"`
}

type Chat struct {
	ChatId   uuid.UUID `json:"chatid"	binding:"required"`
	ChatName string    `json:"chatname"	binding:"required,max=30"`
	Owner    uuid.UUID `json:"id"		binding:"required"`
	Users    []string  `json:"users"`
}

type UserLoginData struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Token struct {
	TokenString string `json:"accessToken"`
}

type UserEmailData struct {
	Email string `json:"email" binding:"required"`
}

type UserChangePassData struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type UserChangeDescriptionData struct {
	Description string `json:"description"`
}

func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{
		pool: pool,
	}
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func comparePasswords(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return nil
	}
	// Если пароли совпадают, возвращаем ошибку.
	return fmt.Errorf("пароли совпадают")
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

func (db DB) RegisterUser(userData User) (string, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return "Проблема с установкой соединения", fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	userData.UserId = uuid.New()
	password, hashErr := hashPassword(userData.Password)
	if hashErr != nil {
		return "Проблема с хешированием", fmt.Errorf("unable to hashPass: %v", hashErr)
	}

	err = conn.QueryRow(context.Background(),
		`INSERT INTO users(userid, username, email, password) VALUES ($1, $2, $3, $4) RETURNING userid`,
		userData.UserId, userData.Username, userData.Email, password).Scan(&userData.UserId)
	if err != nil {
		return "Проблема с запросом в базу данных", fmt.Errorf("unable to INSERT: %v", err)
	}

	return "Вы успешно зарегистрировались", nil
}

func (db DB) GetUserByEmail(email string) (*User, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	var user User
	err = conn.QueryRow(context.Background(), "SELECT userid, username, isadmin, email, password FROM users WHERE email = $1", email).
		Scan(&user.UserId, &user.Username, &user.IsAdmin, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve user: %v", err)
	}

	return &user, err
}

func (db DB) userExists(userID uuid.UUID) (bool, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return false, fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	var exists bool
	err = conn.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT 1 FROM users WHERE userid = $1)", userID).
		Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %v", err)
	}

	return exists, err
}

func (db DB) GetUserInfo(userID uuid.UUID) ([]User, error) {
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
		"SELECT userid, \"username\", \"email\", isadmin, description, avatar FROM users WHERE userid = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from database: %v", err)
	}
	defer rows.Close()

	var data []User
	for rows.Next() {
		var d User
		err = rows.Scan(&d.UserId, &d.Username, &d.Email, &d.IsAdmin, &d.Description, &d.Avatar)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %v", err)
		}
		data = append(data, d)
	}
	return data, err
}

func (db DB) ChangePassword(email string, password string) error {
	exists, err := db.userExistsByEmail(email)
	if err != nil {
		return fmt.Errorf("unable to check if user exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("user with email %s does not exist", email)
	}
	user, err := db.getUserByEmail(email)
	if err != nil {
		return fmt.Errorf("unable to retrieve user data: %v", err)
	}
	er := comparePasswords(user.Password, password)
	if er != nil {
		return fmt.Errorf("%v", er)
	}
	hashPass, er := hashPassword(password)
	if er != nil {
		return fmt.Errorf("%v", err)
	}
	_, err = db.pool.Exec(context.Background(), "UPDATE users SET password = $1 WHERE email = $2", hashPass, email)
	if err != nil {
		return fmt.Errorf("unable to update password: %v", err)
	}
	return nil
}

func (db DB) userExistsByEmail(email string) (bool, error) {
	var exists bool
	err := db.pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db DB) getUserByEmail(email string) (*User, error) {
	var user User
	err := db.pool.QueryRow(context.Background(), "SELECT userid, username, email, password, isadmin, description, avatar FROM users WHERE email = $1", email).Scan(&user.UserId, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.Description, &user.Avatar)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db DB) CreateChat(chat *Chat) error {
	chat.ChatId = uuid.New()

	_, err := db.pool.Exec(context.Background(), "INSERT INTO chat (chatid, chatName, owner) VALUES ($1, $2, $3)", chat.ChatId, chat.ChatName, chat.Owner)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) AddUserToChat(userID string, chatID uuid.UUID) error {
	_, err := db.pool.Exec(context.Background(), "INSERT INTO user_chat (userid, chatid) VALUES ($1, $2)", userID, chatID)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) ChangeDescription(userID string, description string) error {
	_, err := db.pool.Exec(context.Background(), `
		UPDATE users
		SET description = $1
		WHERE userid = $2
	`, description, userID)

	if err != nil {
		return fmt.Errorf("unable to update description: %v", err)
	}
	return nil
}
