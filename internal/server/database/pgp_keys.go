package database

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// CreatePGPKey will create PGP Key
func (db *Database) CreatePGPKey(userID uint, fingerprint, publicKey string) (*models.PGPKey, error) {
	pgpKey := &models.PGPKey{
		UserID:      userID,
		Fingerprint: fingerprint,
		PublicKey:   publicKey,
	}

	if err := db.gormDB.Create(pgpKey).Error; err != nil {
		return nil, err
	}

	return pgpKey, nil
}

// GetAllPGPKeysByUserID returns all PGPKeys that match the userID
func (db *Database) GetAllPGPKeysByUserID(userID uint) ([]*models.PGPKey, error) {
	var pgpKeys []*models.PGPKey

	query := db.gormDB.Model(
		&models.PGPKey{},
	).Where(
		&models.PGPKey{
			UserID: userID,
		},
	)

	if err := query.Find(&pgpKeys).Error; err != nil {
		return nil, err
	}

	return pgpKeys, nil
}
