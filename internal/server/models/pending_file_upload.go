package models

import (
	"time"
)

// PendingFileUpload is a file that has not yet been associated with
// a source package upload
type PendingFileUpload struct {
	ID         uint
	Filename   string
	SHA256Sum  string
	UploadedAt time.Time
	Completed  bool
}
