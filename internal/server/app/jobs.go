package app

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// GetAllJobs returns all jobs
func (app *App) GetAllJobs() ([]*models.Job, error) {
	return app.db.GetAllJobs()
}
