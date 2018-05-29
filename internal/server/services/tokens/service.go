package tokens

import (
	"crypto/rand"
	"encoding/hex"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

//Service manages access tokens
type Service struct {
	db *database.Database
}

//New creates a new tokens service
func New(db *database.Database) *Service {
	service := &Service{
		db: db,
	}
	return service
}

//CreateToken generates an access token
func (service *Service) CreateToken(userID uint, name string) (*models.AccessToken, error) {
	if name == "" {
		return nil, errors.New("the token name cannot be empty")

	}

	b := make([]byte, 20)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	token := hex.EncodeToString(b)

	accessToken, err := service.db.CreateAccessToken(userID, name, token)
	if err != nil {
		return nil, err
	}
	return accessToken, err
}

// RemoveToken removes an access token
func (service *Service) RemoveToken(id uint, userID uint) error {
	return service.db.RemoveAccessToken(id, userID)
}

// GetUserByToken returns the user associated with the given token
func (service *Service) GetUserByToken(token string) (*models.User, error) {
	tokens, err := service.db.GetAllAccessTokensByToken(token)
	if err != nil {
		return nil, err
	}
	if len(tokens) < 1 {
		return nil, nil
	}

	user, err := service.db.GetUser(tokens[0].UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserTokens returns all access tokens associated with a user
func (service *Service) GetUserTokens(userID uint) ([]*models.AccessToken, error) {
	tokens, err := service.db.GetAllAccessTokensByUserID(userID)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}
