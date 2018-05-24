package models

// Upload is a package upload
type Upload struct {
	ID         uint
	UserID     uint
	Source     string
	Version    string
	Maintainer string
	ChangedBy  string
}
