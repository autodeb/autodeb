package jobs_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services/servicestest"
)

func TestSetJobStatusBuildJobCompleted(t *testing.T) {
	servicesTest := servicestest.SetupTest(t)
	jobsService := servicesTest.Services.Jobs()
	artifactsService := servicesTest.Services.Artifacts()

	// Create an upload
	upload, err := servicesTest.DB.CreateUpload(22, "testsource", "testversion", "testmaintainer", "testchangedby", true)
	assert.NoError(t, err)
	assert.NotNil(t, upload)

	// Create a build job for this upload
	job, err := jobsService.CreateBuildJob(upload.ID)
	assert.NoError(t, err)
	assert.NotNil(t, job)

	// There should be only one job associated with the upload
	jobs, err := jobsService.GetAllJobsByUploadID(upload.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(jobs))

	// Upload debs
	_, err = artifactsService.CreateArtifact(
		job.ID, "test1.deb", strings.NewReader("deb content"),
	)
	assert.NoError(t, err)
	_, err = artifactsService.CreateArtifact(
		job.ID, "test2.deb", strings.NewReader("deb content"),
	)
	assert.NoError(t, err)

	// Mark the job as successfull
	err = jobsService.SetJobStatus(job.ID, models.JobStatusSuccess)
	assert.NoError(t, err)

	// There should now be two new autopkgtest jobs associated with the upload
	jobs, err = jobsService.GetAllJobsByUploadID(upload.ID)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(jobs))
	assert.Equal(t, jobs[1].Type, models.JobTypeAutopkgtest)
	assert.Equal(t, jobs[1].UploadID, upload.ID)
	assert.Equal(t, jobs[2].Type, models.JobTypeAutopkgtest)
	assert.Equal(t, jobs[2].UploadID, upload.ID)
}
