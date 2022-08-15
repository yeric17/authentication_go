package models

import (
	"database/sql"
	"errors"
	"fmt"

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
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token"`
}

type UserClaims struct {
	// ID         string `json:"id"`
	// UniqueName string `json:"unique_name"`
	// Name       string `json:"name"`
	// Email      string `json:"email"`
	// Avatar     string `json:"avatar"`
	// Phone      string `json:"phone,omitempty"`
	// Status     int    `json:"status,omitempty"`
	// Token      string `json:"token,omitempty"`
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

func (u *User) Create() error {

	if err := u.Validate(); err != nil {
		return fmt.Errorf("not valid user: %s", err)
	}

	if u.Avatar == "" {
		u.Avatar = utils.GetImageByLetter(u.Name[0:1])
	}

	query := `INSERT INTO Users (user_unique_name, user_name, user_password, user_email, user_avatar, user_phone, user_status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING user_id`

	passByte, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)

	if err != nil {
		return fmt.Errorf("error creating user: %s", err)
	}

	u.Password = string(passByte)

	err = db.QueryRow(query, u.UniqueName, u.Name, u.Password, u.Email, u.Avatar, u.Phone, 2).Scan(&u.ID)

	if err != nil {
		return fmt.Errorf("error get new user: %s", err)
	}

	u.Password = ""

	return nil
}

func (u *User) GetByEmail() error {
	query := `SELECT user_id, user_unique_name, user_name, user_password, user_email, user_avatar, user_phone, user_status
	FROM Users
	WHERE user_email = $1`

	err := db.QueryRow(query, u.Email).Scan(&u.ID, &u.UniqueName, &u.Name, &u.Password, &u.Email, &u.Avatar, &u.Phone, &u.Status)

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
	query := `SELECT r_token_status FROM Refresh_Tokens WHERE r_token_value = $1`

	var status int
	err := db.QueryRow(query, u.RefreshToken).Scan(&status)

	if err != nil {
		return fmt.Errorf("can not refresh token: %s", err)
	}

	if status == 0 {
		return errors.New("the token is not valid")
	}

	return nil
}
