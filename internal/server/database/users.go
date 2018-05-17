package database

import (
	"github.com/jinzhu/gorm"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// CreateUser will create a user
func (db *Database) CreateUser(id uint, username string) (*models.User, error) {
	user := &models.User{
		ID:       id,
		Username: username,
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
