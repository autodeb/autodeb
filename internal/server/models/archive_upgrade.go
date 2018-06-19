package models

// ArchiveUpgrade is an upgrade of all source packages in the archive to a newer upstream version
type ArchiveUpgrade struct {
	ID           uint `json:"id"`
	UserID       uint `json:"user_id"`
	PackageCount int  `json:"package_count"`
}
