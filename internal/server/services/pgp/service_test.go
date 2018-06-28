package pgp

import (
	"net/url"
	"strings"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/pgp"
	"salsa.debian.org/autodeb-team/autodeb/internal/pgp/pgptest"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database/databasetest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) *Service {
	db := databasetest.SetupTest(t)

	serverURL, err := url.Parse("http://pgptests:8080")
	require.NoError(t, err)

	service := New(db, serverURL)
	return service
}

func TestAddPGPKey(t *testing.T) {
	service := setupTest(t)

	user, err := service.db.CreateUser("testUser", 33)
	assert.NoError(t, err)

	// Verify that the keyring is empty
	keyring, err := service.keyRing()
	assert.Equal(t, 0, len(keyring))
	assert.NoError(t, err)

	// Verify that the key does not yet exist
	keys, err := service.GetUserPGPKeys(user.ID)
	assert.NotNil(t, keys)
	assert.Equal(t, 0, len(keys))
	assert.NoError(t, err)

	// Sign the proof
	proof, err := pgp.Clearsign(
		strings.NewReader(service.ExpectedPGPKeyProofText(user.ID)),
		strings.NewReader(pgptest.TestKeyPrivate),
	)
	assert.NoError(t, err)

	// Add the key
	err = service.AddUserPGPKey(user.ID, pgptest.TestKeyPublic, proof)
	assert.NoError(t, err)

	// Verify that the key was added to the database
	keys, err = service.GetUserPGPKeys(user.ID)
	assert.NotNil(t, keys)
	assert.Equal(t, 1, len(keys))
	assert.NoError(t, err)
	key := keys[0]
	assert.Equal(t, pgptest.TestKeyFingerprint, key.Fingerprint)
	assert.Equal(t, pgptest.TestKeyPublic, key.PublicKey)

	// Verify the keyring's content
	keyring, err = service.keyRing()
	assert.Equal(t, 1, len(keyring))

	// Identify the signer of the proof
	signerID, err := service.IdentifySigner(strings.NewReader(proof))
	assert.NoError(t, err)
	assert.Equal(t, user.ID, signerID)
}

func TestAddPGPKeyAlreadyRegistered(t *testing.T) {
	service := setupTest(t)

	user, err := service.db.CreateUser("testUser", 33)
	assert.NoError(t, err)

	// Sign the proof
	proof, err := pgp.Clearsign(
		strings.NewReader(service.ExpectedPGPKeyProofText(user.ID)),
		strings.NewReader(pgptest.TestKeyPrivate),
	)
	assert.NoError(t, err)

	// Add the key
	err = service.AddUserPGPKey(user.ID, pgptest.TestKeyPublic, proof)
	assert.NoError(t, err)

	// Add the key again
	err = service.AddUserPGPKey(user.ID, pgptest.TestKeyPublic, proof)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is already registered to user")
}

func TestKeyRingEmpty(t *testing.T) {
	service := setupTest(t)
	keyring, err := service.keyRing()
	assert.NoError(t, err, "an empty keyring should not result in an error")
	assert.Nil(t, keyring)
}
