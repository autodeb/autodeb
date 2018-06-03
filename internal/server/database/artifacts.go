package database

import (
	"github.com/jinzhu/gorm"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// CreateArtifact will create a job artifact
func (db *Database) CreateArtifact(jobID uint, filename string) (*models.Artifact, error) {
	artifact := &models.Artifact{
		JobID:    jobID,
		Filename: filename,
	}

	if err := db.gormDB.Create(artifact).Error; err != nil {
		return nil, err
	}

	return artifact, nil
}

// GetArtifact returns the Artifact with the given id
func (db *Database) GetArtifact(id uint) (*models.Artifact, error) {
	var artifact models.Artifact

	query := db.gormDB.Where(
		&models.Artifact{
			ID: id,
		},
	)

	err := query.First(&artifact).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &artifact, nil
}

// GetAllArtifactsByJobID returns all job artifacts for a job
func (db *Database) GetAllArtifactsByJobID(jobID uint) ([]*models.Artifact, error) {
	var artifacts []*models.Artifact

	query := db.gormDB.Model(
		&models.Artifact{},
	).Where(
		&models.Artifact{
			JobID: jobID,
		},
	)

	if err := query.Find(&artifacts).Error; err != nil {
		return nil, err
	}

	return artifacts, nil
}

// GetAllArtifactsByJobIDFilename returns all job artifacts for a job with a matching file name
func (db *Database) GetAllArtifactsByJobIDFilename(jobID uint, filename string) ([]*models.Artifact, error) {
	var artifacts []*models.Artifact

	query := db.gormDB.Model(
		&models.Artifact{},
	).Where(
		&models.Artifact{
			JobID:    jobID,
			Filename: filename,
		},
	)

	if err := query.Find(&artifacts).Error; err != nil {
		return nil, err
	}

	return artifacts, nil
}
