package models

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/yeric17/thullo/pkg/data"
	"github.com/yeric17/thullo/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	db *sql.DB
)

func init() {
	db = data.GetConnection()
}

type User struct {
	ID           string `json:"id"`
	UniqueName   string `json:"unique_name"`
	Name         string `json:"name"`
	Password     string `json:"password,omitempty"`
	Email        string `json:"email"`
	Avatar       string `json:"avatar"`
	Phone        string `json:"phone,omitempty"`
	Status       int    `json:"status,omitempty"`
	Provider     string `json:"provider,omitempty"`
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token"`
}

type UserClaims struct {
	DataUser User `json:"user"`
	jwt.StandardClaims
}

type UserStatus int

const (
	UserStatusInactive UserStatus = iota
	UserStatusActive
	UserStatusWaitingConfirmation
)

func (u *User) Validate() error {
	if u.UniqueName == "" {
		return fmt.Errorf("unique_name is required")
	}
	if u.Name == "" {
		return fmt.Errorf("name is required")
	}
	if u.Password == "" {
		return fmt.Errorf("password is required")
	}
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	return nil
}

func (u *User) Create(provider string) error {

	if err := u.Validate(); err != nil {
		return fmt.Errorf("not valid user: %s", err)
	}

	if u.Avatar == "" {
		u.Avatar = utils.GetImageByLetter(u.Name[0:1])
	}

	query := `INSERT INTO Users (user_unique_name, user_name, user_password, user_email, user_avatar, user_phone, user_status, user_auth_provider) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING user_id`

	passByte, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)

	if err != nil {
		return fmt.Errorf("error creating user: %s", err)
	}

	u.Password = string(passByte)

	var status int
	if provider == "google" {
		status = 1
	} else {
		status = 2
	}

	err = db.QueryRow(query, u.UniqueName, u.Name, u.Password, u.Email, u.Avatar, u.Phone, status, provider).Scan(&u.ID)

	if err != nil {
		return fmt.Errorf("error get new user: %s", err)
	}

	u.Password = ""

	return nil
}

func (u *User) GetByEmail(withPass bool) error {
	var query string
	var err error
	if withPass {
		query = `SELECT user_id, user_unique_name, user_name, user_password, user_email, user_avatar, user_phone, user_status
		FROM Users
		WHERE user_email = $1`
		err = db.QueryRow(query, u.Email).Scan(&u.ID, &u.UniqueName, &u.Name, &u.Password, &u.Email, &u.Avatar, &u.Phone, &u.Status)
	} else {
		query = `SELECT user_id, user_unique_name, user_name, user_email, user_avatar, user_phone, user_status
		FROM Users
		WHERE user_email = $1`
		err = db.QueryRow(query, u.Email).Scan(&u.ID, &u.UniqueName, &u.Name, &u.Email, &u.Avatar, &u.Phone, &u.Status)
	}

	if err != nil {
		return fmt.Errorf("error getting user: %s", err)
	}

	return nil
}

func (u *User) GetByID() error {

	query := `SELECT user_id, user_unique_name, user_name, user_email, user_avatar, user_phone, user_status
	FROM Users
	WHERE user_id = $1`

	err := db.QueryRow(query, u.ID).Scan(&u.ID, &u.UniqueName, &u.Name, &u.Email, &u.Avatar, &u.Phone, &u.Status)

	if err != nil {
		return fmt.Errorf("error getting user: %s", err)
	}

	return nil
}

func (u *User) CreateRefreshToken() error {
	query := `UPDATE Refresh_Tokens SET r_token_status = 0 WHERE r_token_user_id = $1`

	_, err1 := db.Exec(query, u.ID)

	if err1 != nil {
		return fmt.Errorf("can not disbled refresh token: %s", err1)
	}

	queryInsert := `INSERT INTO Refresh_Tokens (r_token_user_id, r_token_status) 
	VALUES ($1, $2) 
	RETURNING r_token_value`

	fmt.Println(u.ID)
	err := db.QueryRow(queryInsert, u.ID, 1).Scan(&u.RefreshToken)

	if err != nil {
		return fmt.Errorf("can not create refresh token: %s", err)
	}
	return nil
}

func (u *User) ValidRefreshToken() error {
	query := `SELECT r_token_status, r_token_user_id FROM Refresh_Tokens WHERE r_token_value = $1`

	err := db.QueryRow(query, u.RefreshToken).Scan(&u.Status, &u.ID)

	if err != nil {
		return fmt.Errorf("can not refresh token: %s", err)
	}

	if u.Status == 0 {
		return errors.New("the token is not valid")
	}

	if err := u.GetByID(); err != nil {
		return fmt.Errorf("user does not exists: %s", err)
	}

	return nil
}

func (u *User) ValidToken() error {
	var claim UserClaims

	_, err := jwt.ParseWithClaims(u.Token, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return fmt.Errorf("error authenticating: %s", err)
	}

	return nil
}
