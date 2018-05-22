package models

// PGPKey represents the pgp key of a user
type PGPKey struct {
	ID          uint
	UserID      uint
	Fingerprint string
	PublicKey   string
}
