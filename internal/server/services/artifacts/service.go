package artifacts

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

//Service manages job artifacts
type Service struct {
	db *database.Database
	fs filesystem.FS
}

//New creates an artifacts service
func New(db *database.Database, fs filesystem.FS) *Service {
	service := &Service{
		db: db,
		fs: fs,
	}
	return service
}

func (service *Service) artifactsPath() string {
	return "/"
}

func (service *Service) artifactPath(artifactID uint) string {
	artifactPath := filepath.Join(
		service.artifactsPath(),
		fmt.Sprint(artifactID),
	)
	return artifactPath
}

// CreateArtifact creates a new artifact
func (service *Service) CreateArtifact(jobID uint, filename string, content io.Reader) (*models.Artifact, error) {
	artifact, err := service.db.CreateArtifact(jobID, filename)
	if err != nil {
		return nil, err
	}

	f, err := service.fs.Create(
		service.artifactPath(artifact.ID),
	)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := io.Copy(f, content); err != nil {
		return nil, err
	}

	return artifact, nil
}

// GetArtifact returns an artifact by id
func (service *Service) GetArtifact(id uint) (*models.Artifact, error) {
	artifact, err := service.db.GetArtifact(id)
	if err != nil {
		return nil, err
	}
	return artifact, nil
}

// GetArtifactContent returns the content of an artifact
func (service *Service) GetArtifactContent(artifactID uint) (io.ReadCloser, error) {
	file, err := service.fs.Open(service.artifactPath(artifactID))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return file, nil
}

// GetAllArtifactsByJobIDFilename returns artifacts matching the job id and file name
func (service *Service) GetAllArtifactsByJobIDFilename(jobID uint, filename string) ([]*models.Artifact, error) {
	return service.db.GetAllArtifactsByJobIDFilename(jobID, filename)
}

// GetAllArtifactsByJobID returns all artifacts by job id
func (service *Service) GetAllArtifactsByJobID(jobID uint) ([]*models.Artifact, error) {
	artifacts, err := service.db.GetAllArtifactsByJobID(jobID)
	if err != nil {
		return nil, err
	}
	return artifacts, nil
}
