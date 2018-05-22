package servicestest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services"
)

//ServicesTest makes it easy to test services
type ServicesTest struct {
	DataFS    filesystem.FS
	DB        *database.Database
	ServerURL string
	Services  *services.Services
}

//SetupTest will create a test App
func SetupTest(t *testing.T) *ServicesTest {
	serverURL := "https://test.auto.debian.net"

	db, err := database.NewDatabase("sqlite3", ":memory:")
	require.NoError(t, err)

	dataFS := filesystem.NewMemMapFS()

	services, err := services.New(
		db,
		dataFS,
		serverURL,
	)
	require.NoError(t, err)

	appTest := &ServicesTest{
		DataFS:    dataFS,
		DB:        db,
		ServerURL: serverURL,
		Services:  services,
	}

	return appTest
}
