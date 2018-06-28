package models

// User is a user of the service
type User struct {
	ID                uint   `json:"id"`
	Username          string `json:"username"`
	AuthBackendUserID uint   `json:"auth_backend_user_id"`
}
