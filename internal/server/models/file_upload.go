package models

import (
	"time"
)

// FileUpload is an individual file that was uploaded
type FileUpload struct {
	ID         uint      `json:"id"`
	Filename   string    `json:"filename"`
	SHA256Sum  string    `json:"sha256sum"`
	UploadedAt time.Time `json:"uploaded_at"`
	Completed  bool      `json:"completed"`
	UploadID   uint      `json:"upload_id"`
}
