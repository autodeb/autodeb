package database

import (
	"github.com/jinzhu/gorm"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// CreateUser will create a user
func (db *Database) CreateUser(username string, authBackendUserID uint) (*models.User, error) {
	user := &models.User{
		Username:          username,
		AuthBackendUserID: authBackendUserID,
	}

	if err := db.gormDB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser returns the User with the given id
func (db *Database) GetUser(id uint) (*models.User, error) {
	var user models.User

	query := db.gormDB.Where(
		&models.User{
			ID: id,
		},
	)

	err := query.First(&user).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByAuthBackendUserID returns the User with a matching AuthBackendUserID
func (db *Database) GetUserByAuthBackendUserID(id uint) (*models.User, error) {
	var user models.User

	query := db.gormDB.Where(
		&models.User{
			AuthBackendUserID: id,
		},
	)

	err := query.First(&user).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
