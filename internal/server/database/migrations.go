package database

import (
	"github.com/jinzhu/gorm"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func runMigrations(gormDB *gorm.DB) {
	gormDB.AutoMigrate(&models.AccessToken{})
	gormDB.AutoMigrate(&models.ArchiveUpgrade{})
	gormDB.AutoMigrate(&models.Artifact{})
	gormDB.AutoMigrate(&models.FileUpload{})
	gormDB.AutoMigrate(&models.Job{})
	gormDB.AutoMigrate(&models.PGPKey{})
	gormDB.AutoMigrate(&models.Upload{})
	gormDB.AutoMigrate(&models.User{})
}
