package servicestest

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/pgp/pgptest"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database/databasetest"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services"
)

//ServicesTest makes it easy to test services
type ServicesTest struct {
	t           *testing.T
	DataFS      filesystem.FS
	DB          *database.Database
	ServerURL   *url.URL
	AptlyAPIURL *url.URL
	Services    *services.Services
}

//SetupTest will create a test App
func SetupTest(t *testing.T) *ServicesTest {
	serverURL, err := url.Parse("https://test.auto.debian.net")
	require.NoError(t, err)

	aptlyAPIURL, err := url.Parse("https://aptly.test")
	require.NoError(t, err)

	db := databasetest.SetupTest(t)

	dataFS := filesystem.NewMemMapFS()

	services, err := services.New(
		db,
		dataFS,
		serverURL,
	)
	require.NoError(t, err)

	appTest := &ServicesTest{
		t:           t,
		DataFS:      dataFS,
		DB:          db,
		ServerURL:   serverURL,
		AptlyAPIURL: aptlyAPIURL,
		Services:    services,
	}

	return appTest
}

// GetOrCreateTestUser will get or create the test user
// and return the user, its pgp key and its private key
func (servicesTest *ServicesTest) GetOrCreateTestUser() *models.User {
	user, err := servicesTest.DB.GetUser(uint(1))
	require.NoError(servicesTest.t, err)

	if user != nil {
		return user
	}

	user, err = servicesTest.DB.CreateUser("testuser3579", 13)
	require.NoError(servicesTest.t, err)
	require.NotNil(servicesTest.t, user)

	return user
}

// AddTokenToUser will generate an access token for the given user
func (servicesTest *ServicesTest) AddTokenToUser(user *models.User) *models.AccessToken {
	token, err := servicesTest.Services.Tokens().CreateToken(user.ID, "test token")
	require.NoError(servicesTest.t, err)
	return token
}

// AddPGPKeyToUser will add a pgp key to the user and return public and private key
func (servicesTest *ServicesTest) AddPGPKeyToUser(user *models.User) (*models.PGPKey, string) {
	key, err := servicesTest.DB.CreatePGPKey(
		user.ID,
		pgptest.TestKeyFingerprint,
		pgptest.TestKeyPublic,
	)
	require.NoError(servicesTest.t, err)
	require.NotNil(servicesTest.t, key)

	return key, pgptest.TestKeyPrivate
}
