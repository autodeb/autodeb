package jobs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services/servicestest"
)

func TestCreateArchiveUpgrade(t *testing.T) {
	servicesTest := servicestest.SetupTest(t)
	jobsService := servicesTest.Services.Jobs()

	jobs, err := jobsService.GetAllJobs()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(jobs))

	archiveBackport, err := jobsService.CreateArchiveUpgrade(33, 100)
	assert.NoError(t, err)

	jobs, err = jobsService.GetAllJobs()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(jobs))

	job := jobs[0]
	assert.Equal(t, models.JobParentTypeArchiveUpgrade, job.ParentType)
	assert.Equal(t, archiveBackport.ID, job.ParentID)
}
