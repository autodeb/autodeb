package uploads

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
)

//Manager handles uploads
type Manager struct {
	db     *database.Database
	dataFS filesystem.FS
}

//NewManager creates a new upload manager
func NewManager(db *database.Database, dataFS filesystem.FS) *Manager {
	man := &Manager{
		db:     db,
		dataFS: dataFS,
	}
	return man
}

// UploadedFilesDirectory contains files that are not yet associated
// with a package upload.
func (man *Manager) UploadedFilesDirectory() string {
	return "/files"
}

// UploadsDirectory contains completed uploads.
func (man *Manager) UploadsDirectory() string {
	return "/uploads"
}
