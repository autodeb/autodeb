package models

// ArchiveBackport is a backport of all source packages in the archive to stable
type ArchiveBackport struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
}
