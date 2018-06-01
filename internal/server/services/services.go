//Package services implements the bulk of the application logic.
package services

import (
	"net/url"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services/jobs"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services/pgp"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services/tokens"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services/uploads"
)

//Services holds the services
type Services struct {
	jobs    *jobs.Service
	pgp     *pgp.Service
	uploads *uploads.Service
	tokens  *tokens.Service
}

// New returns a new set of services
func New(db *database.Database, dataFS filesystem.FS, serverURL *url.URL) (*Services, error) {
	// Tokens
	tokensService := tokens.New(db)

	// PGP
	pgpService := pgp.New(db, serverURL)

	// Jobs
	if err := dataFS.MkdirAll("jobs", 0744); err != nil {
		return nil, errors.WithMessage(err, "could not create jobs folder")
	}
	jobsService := jobs.New(
		db,
		filesystem.NewBasePathFS(dataFS, "jobs"),
	)

	// Uploads
	if err := dataFS.MkdirAll("uploads", 0744); err != nil {
		return nil, errors.WithMessage(err, "could not create uploads folder")
	}
	uploadsService := uploads.New(
		db,
		pgpService,
		jobsService,
		filesystem.NewBasePathFS(dataFS, "uploads"),
	)

	services := &Services{
		tokens:  tokensService,
		jobs:    jobsService,
		pgp:     pgpService,
		uploads: uploadsService,
	}

	return services, nil
}

// Tokens returns the tokens service
func (services *Services) Tokens() *tokens.Service {
	return services.tokens
}

// PGP returns the pgp service
func (services *Services) PGP() *pgp.Service {
	return services.pgp
}

// Jobs returns the jobs service
func (services *Services) Jobs() *jobs.Service {
	return services.jobs
}

// Uploads returns the uploads service
func (services *Services) Uploads() *uploads.Service {
	return services.uploads
}
