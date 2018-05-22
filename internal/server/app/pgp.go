package app

import (
	"fmt"
	"strings"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/pgp"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// ExpectedPGPKeyProofText returns the expected PGP Key ownership proof text.
// for a given user.
func (app *App) ExpectedPGPKeyProofText(userID uint) string {
	expectedProofText := fmt.Sprintf(
		"I am User ID %d on %s",
		userID,
		app.config.ServerURL,
	)
	return expectedProofText
}

// AddUserKey associates a PGP key with the user, if the proof is valid.
func (app *App) AddUserKey(userID uint, key, proof string) error {
	signedProofText, entity, err := pgp.VerifySignatureClearsigned(
		strings.NewReader(proof),
		strings.NewReader(key),
	)
	if err != nil {
		return errors.WithMessage(err, "couldn't verify signature")
	}

	signedProofText = strings.TrimSpace(signedProofText)

	if signedProofText != app.ExpectedPGPKeyProofText(userID) {
		return errors.Errorf("Signed proof text did not match the expected proof text")
	}

	fingerprint := pgp.EntityFingerprint(entity)

	if _, err := app.db.CreatePGPKey(userID, fingerprint); err != nil {
		return err
	}

	return nil
}

// GetUserPGPKeys returns all PGP Keys associated with a user
func (app *App) GetUserPGPKeys(userID uint) ([]*models.PGPKey, error) {
	keys, err := app.db.GetAllPGPKeysByUserID(userID)
	if err != nil {
		return nil, err
	}
	return keys, nil
}
