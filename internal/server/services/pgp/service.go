package pgp

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
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
	serverURL *url.URL
}

//New creates a new pgp service
func New(db *database.Database, serverURL *url.URL) *Service {
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
		service.serverURL.String(),
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

	// Ensure that it is not already registered
	if keys, err := service.db.GetAllPGPKeysByFingerprint(fingerprint); err != nil {
		return err
	} else if len(keys) > 0 {
		return errors.Errorf("key %s is already registered to user %d", fingerprint, keys[0].UserID)
	}

	// Add the key to the database
	if _, err := service.db.CreatePGPKey(userID, fingerprint, key); err != nil {
		return err
	}

	return nil
}

// IdentifySigner will identify the user that has signed the given clearsigned
// message.
func (service *Service) IdentifySigner(message io.Reader) (uint, error) {
	keyRing, err := service.keyRing()
	if err != nil {
		return 0, err
	}

	_, entity, err := pgp.VerifySignatureClearsignedKeyRing(message, keyRing)
	if err != nil {
		return 0, err
	}

	keys, err := service.db.GetAllPGPKeysByFingerprint(
		pgp.EntityFingerprint(entity),
	)
	if err != nil {
		return 0, err
	}
	if len(keys) == 0 {
		return 0, errors.New("could not find key matching the signer")
	}

	return keys[0].UserID, nil
}

// keyRing returns keyring that contains all known public keys
func (service *Service) keyRing() (pgp.EntityList, error) {
	pgpKeys, err := service.db.GetAllPGPKeys()
	if err != nil {
		return nil, err
	}

	if len(pgpKeys) == 0 {
		var keyring []*pgp.Entity
		return keyring, nil
	}

	// Read all keys and serialize them to a buffer
	var buf bytes.Buffer
	for _, pgpKey := range pgpKeys {
		keyRing, err := pgp.ReadArmoredKeyRing(
			strings.NewReader(pgpKey.PublicKey),
		)
		if err != nil {
			return nil, errors.WithMessagef(err, "could not load key id %d with fingerprint %s", pgpKey.ID, pgpKey.Fingerprint)
		}
		if err := keyRing[0].Serialize(&buf); err != nil {
			return nil, errors.WithMessagef(err, "could not serialize key id %d with fingerprint %s", pgpKey.ID, pgpKey.Fingerprint)
		}
	}

	// Read the buffer as one keyRing
	keyring, err := pgp.ReadKeyRing(&buf)
	if err != nil {
		return nil, err
	}

	return keyring, err
}

// RemovePGPKey removes the PGP Key with a matching userID and ID
func (service *Service) RemovePGPKey(id uint, userID uint) error {
	return service.db.RemovePGPKey(id, userID)
}

// GetUserPGPKeys returns all PGP Keys associated with a user
func (service *Service) GetUserPGPKeys(userID uint) ([]*models.PGPKey, error) {
	keys, err := service.db.GetAllPGPKeysByUserID(userID)
	if err != nil {
		return nil, err
	}
	return keys, nil
}
