package models

// AccessToken allows API access
type AccessToken struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
	Name   string `json:"name"`
}
