package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/database/databasetest"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func TestChangeJobStatus(t *testing.T) {
	dbTest := databasetest.SetupTest(t)

	job, err := dbTest.CreateJob(models.JobTypeBuild, "", models.JobParentTypeUpload, 22)
	assert.NoError(t, err)
	assert.NotNil(t, job)
	assert.Equal(t, models.JobStatusQueued, job.Status)

	// Change the status once
	err = dbTest.ChangeJobStatus(job.ID, models.JobStatusAssigned)
	assert.NoError(t, err)

	// Get the job and see if the status was updated
	job, err = dbTest.GetJob(job.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.JobStatusAssigned, job.Status)

	// Change the status again, should fail
	err = dbTest.ChangeJobStatus(job.ID, models.JobStatusAssigned)
	assert.Error(t, err)
}
