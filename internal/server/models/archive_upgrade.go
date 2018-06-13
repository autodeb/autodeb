package models

// ArchiveUpgrade is an upgrade of all source packages in the archive to a newer upstream version
type ArchiveUpgrade struct {
	ID           uint `json:"id"`
	UserID       uint `json:"user_id"`
	PackageCount uint `json:"package_count"`
}
