package pgp

import (
	"fmt"
	"strings"
	"time"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/pgp"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

//Service manages pgp keys and verification
type Service struct {
	db        *database.Database
	serverURL string
}

//New creates a new pgp service
func New(db *database.Database, serverURL string) *Service {
	service := &Service{
		db:        db,
		serverURL: serverURL,
	}
	return service
}

// ExpectedPGPKeyProofText returns the expected PGP Key ownership proof text.
// for a given user.
func (service *Service) ExpectedPGPKeyProofText(userID uint) string {
	expectedProofText := fmt.Sprintf(
		"As of %s, I am User ID %d on %s",
		time.Now().Format("2006-01-02"),
		userID,
		service.serverURL,
	)
	return expectedProofText
}

// AddUserPGPKey associates a PGP key with the user, if the proof is valid.
func (service *Service) AddUserPGPKey(userID uint, key, proof string) error {

	// Read the provided keyring
	keyring, err := pgp.ReadArmoredKeyRing(strings.NewReader(key))
	if err != nil {
		return errors.WithMessage(err, "could not read the provided key")
	}

	// The keyring should only contain one key
	if numKeys := len(keyring); numKeys != 1 {
		return errors.Errorf("expected 1 key, %d were provided", numKeys)
	}

	// The key can have one self-sig, nothing more.
	if numSignatures := len(pgp.EntitySignatures(keyring[0])); numSignatures > 1 {
		return errors.Errorf("the provided key should be minimal but it has %d signatures on it", numSignatures)
	}

	// Verify the signature on the proof
	signedProofText, entity, err := pgp.VerifySignatureClearsignedKeyRing(
		strings.NewReader(proof),
		keyring,
	)
	if err != nil {
		return errors.WithMessage(err, "couldn't verify signature")
	}

	// Verify that the signed proof matches the expected proof text
	signedProofText = strings.TrimSpace(signedProofText)
	if signedProofText != service.ExpectedPGPKeyProofText(userID) {
		return errors.Errorf("Signed proof text did not match the expected proof text")
	}

	// Get the fingerpring of the key
	fingerprint := pgp.EntityFingerprint(entity)

	// Add the key to the database
	if _, err := service.db.CreatePGPKey(userID, fingerprint, key); err != nil {
		return err
	}

	return nil
}

// GetUserPGPKeys returns all PGP Keys associated with a user
func (service *Service) GetUserPGPKeys(userID uint) ([]*models.PGPKey, error) {
	keys, err := service.db.GetAllPGPKeysByUserID(userID)
	if err != nil {
		return nil, err
	}
	return keys, nil
}
