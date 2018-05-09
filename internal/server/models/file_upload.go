package models

import (
	"time"
)

// FileUpload is an individual file that was uploaded
type FileUpload struct {
	ID         uint
	Filename   string
	SHA256Sum  string
	UploadedAt time.Time
	Completed  bool
	UploadID   uint
}
