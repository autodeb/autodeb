package database

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// CreateAccessToken will create an access token
func (db *Database) CreateAccessToken(userID uint, name, token string) (*models.AccessToken, error) {
	accessToken := &models.AccessToken{
		UserID: userID,
		Name:   name,
		Token:  token,
	}

	if err := db.gormDB.Create(accessToken).Error; err != nil {
		return nil, err
	}

	return accessToken, nil
}

// RemoveAccessToken removes all matching access tokens
func (db *Database) RemoveAccessToken(id, userID uint) error {
	query := db.gormDB.Model(
		&models.AccessToken{},
	).Where(
		&models.PGPKey{
			ID:     id,
			UserID: userID,
		},
	)

	if err := query.Delete(&models.AccessToken{}).Error; err != nil {
		return err
	}

	return nil
}

// GetAllAccessTokensByUserID returns all AccessTokens that match the userID
func (db *Database) GetAllAccessTokensByUserID(userID uint) ([]*models.AccessToken, error) {
	var accessTokens []*models.AccessToken

	query := db.gormDB.Model(
		&models.AccessToken{},
	).Where(
		&models.AccessToken{
			UserID: userID,
		},
	)

	if err := query.Find(&accessTokens).Error; err != nil {
		return nil, err
	}

	return accessTokens, nil
}

// GetAllAccessTokensByToken returns all AccessTokens that match the userID
func (db *Database) GetAllAccessTokensByToken(token string) ([]*models.AccessToken, error) {
	var accessTokens []*models.AccessToken

	query := db.gormDB.Model(
		&models.AccessToken{},
	).Where(
		&models.AccessToken{
			Token: token,
		},
	)

	if err := query.Find(&accessTokens).Error; err != nil {
		return nil, err
	}

	return accessTokens, nil
}
