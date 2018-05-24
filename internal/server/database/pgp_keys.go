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

// GetAllPGPKeys returns all PGP Keys
func (db *Database) GetAllPGPKeys() ([]*models.PGPKey, error) {
	var pgpKeys []*models.PGPKey

	if err := db.gormDB.Model(&models.PGPKey{}).Find(&pgpKeys).Error; err != nil {
		return nil, err
	}

	return pgpKeys, nil
}

// GetAllPGPKeysByFingerprint returns all keys that match the given fingerprint
func (db *Database) GetAllPGPKeysByFingerprint(fingerprint string) ([]*models.PGPKey, error) {
	var pgpKeys []*models.PGPKey

	query := db.gormDB.Model(
		&models.PGPKey{},
	).Where(
		&models.PGPKey{
			Fingerprint: fingerprint,
		},
	)

	if err := query.Find(&pgpKeys).Error; err != nil {
		return nil, err
	}

	return pgpKeys, nil
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
