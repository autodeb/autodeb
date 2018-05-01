package app

import (
	"testing"

	"github.com/stretchr/testify/require"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
)

func setupTest(t *testing.T) (*App, filesystem.FS, *database.Database) {
	dataFS := filesystem.NewMemMapFs()

	db, err := database.NewDatabase(
		&database.Config{
			Driver:           "sqlite3",
			ConnectionString: ":memory:",
		},
	)
	require.NoError(t, err)

	app, err := NewApp(db, dataFS)
	require.NoError(t, err)

	return app, dataFS, db
}
