// Package database is responsible for all database logic. It handles
// connecting to the database and performing queries.
package database

import (
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite" // sqlite database driver
)

// Database is an implementation of the app.Datastore interface
type Database struct {
	gormDB *gorm.DB
}

// NewDatabase creates a new database connection
func NewDatabase(driver, connectionString string) (*Database, error) {

	gormDB, err := gorm.Open(driver, connectionString)
	if err != nil {
		return nil, err
	}

	runMigrations(gormDB)

	database := Database{
		gormDB: gormDB,
	}

	return &database, nil
}
