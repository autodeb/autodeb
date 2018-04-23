package database

import (
	"github.com/jinzhu/gorm"

	"salsa.debian.org/aviau/autodeb/internal/server/models"
)

func runMigrations(gormDB *gorm.DB) {
	gormDB.AutoMigrate(&models.Upload{})
}
