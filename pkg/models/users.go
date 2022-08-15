package models

import (
	"database/sql"
	"fmt"

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
	ID         string `json:"id"`
	UniqueName string `json:"unique_name"`
	Name       string `json:"name"`
	Password   string `json:"password,omitempty"`
	Email      string `json:"email"`
	Avatar     string `json:"avatar"`
	Phone      string `json:"phone,omitempty"`
}

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
