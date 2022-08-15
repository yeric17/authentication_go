package models

type RefreshToken struct {
	UserID string `json:"user_id"`
	Value  string `json:"value"`
}
