package uploads

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services/pgp"
)

//Service manages uploads
type Service struct {
	db         *database.Database
	pgpService *pgp.Service
	fs         filesystem.FS
}

//New creates a new upload service
func New(db *database.Database, pgpService *pgp.Service, fs filesystem.FS) *Service {
	service := &Service{
		db:         db,
		pgpService: pgpService,
		fs:         fs,
	}
	return service
}

// FS returns the services's filesystem
func (service *Service) FS() filesystem.FS {
	return service.fs
}

// UploadedFilesDirectory contains files that are not yet associated
// with a package upload.
func (service *Service) UploadedFilesDirectory() string {
	return "/files"
}

// UploadsDirectory contains completed uploads.
func (service *Service) UploadsDirectory() string {
	return "/uploads"
}

// GetAllUploads returns all uploads
func (service *Service) GetAllUploads() ([]*models.Upload, error) {
	return service.db.GetAllUploads()
}

//GetAllFileUploadsByUploadID returns all FileUploads associated to an upload
func (service *Service) GetAllFileUploadsByUploadID(uploadID uint) ([]*models.FileUpload, error) {
	fileUploads, err := service.db.GetAllFileUploadsByUploadID(uploadID)
	if err != nil {
		return nil, err
	}
	return fileUploads, nil
}

// GetUploadDSC returns the DSC of the upload with a matching id
func (service *Service) GetUploadDSC(uploadID uint) (io.ReadCloser, error) {
	fileUploads, err := service.GetAllFileUploadsByUploadID(uploadID)
	if err != nil {
		return nil, err
	}

	for _, fileUpload := range fileUploads {
		if strings.HasSuffix(fileUpload.Filename, ".dsc") {
			return service.GetUploadFile(uploadID, fileUpload.Filename)
		}
	}

	return nil, nil
}

// GetUploadFile returns the file associated with the upload id and filename
func (service *Service) GetUploadFile(uploadID uint, filename string) (io.ReadCloser, error) {
	file, err := service.fs.Open(
		filepath.Join(
			service.UploadsDirectory(),
			fmt.Sprint(uploadID),
			filename,
		),
	)
	if err != nil {
		return nil, err
	}
	return file, nil
}
