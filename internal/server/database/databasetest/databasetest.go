package databasetest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
)

// SetupTest will create a new database based on an in-memory backend for
// testing
func SetupTest(t *testing.T) *database.Database {
	db, err := database.NewDatabase("sqlite3", ":memory:")
	require.NoError(t, err)
	return db
}
