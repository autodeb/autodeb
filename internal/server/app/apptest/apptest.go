package apptest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
)

//SetupTest will create a test App
func SetupTest(t *testing.T) (*app.App, filesystem.FS, *database.Database) {
	dataFS := filesystem.NewMemMapFs()

	db, err := database.NewDatabase(
		"sqlite3",
		":memory:",
	)
	require.NoError(t, err)

	app, err := app.NewApp(db, dataFS)
	require.NoError(t, err)

	return app, dataFS, db
}
