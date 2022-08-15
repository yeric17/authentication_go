package models

type User struct {
	ID         string `json:"id"`
	UniqueName string `json:"unique_name"`
	Name       string `json:"name"`
	Password   string `json:"password,omitempty"`
	Email      string `json:"email"`
	Avatar     string `json:"avatar"`
	Phone      string `json:"phone,omitempty"`
}
