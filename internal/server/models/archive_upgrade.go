package models

import (
	"fmt"
)

// ArchiveUpgrade is an upgrade of all source packages in the archive to a newer upstream version
type ArchiveUpgrade struct {
	ID           uint `json:"id"`
	UserID       uint `json:"user_id"`
	PackageCount int  `json:"package_count"`
}

// RepositoryName is the ArchiveUpgrade's Aptly repository name
func (archiveUpgrade *ArchiveUpgrade) RepositoryName() string {
	return fmt.Sprintf("archive-upgrade-%d", archiveUpgrade.ID)
}
