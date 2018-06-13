package jobs_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services/servicestest"
)

func TestProcessJobStatusBuildAndDontForward(t *testing.T) {
	servicesTest := servicestest.SetupTest(t)
	jobsService := servicesTest.Services.Jobs()

	// Create an upload
	upload, err := servicesTest.DB.CreateUpload(22, "testsource", "testversion", "testmaintainer", "testchangedby", false, true)
	assert.NoError(t, err)
	assert.NotNil(t, upload)

	// Create a build job for this upload
	job, err := jobsService.CreateJob(
		models.JobTypeBuild, "", models.JobParentTypeUpload, upload.ID,
	)
	assert.NoError(t, err)
	assert.NotNil(t, job)

	// There should be only one job associated with the upload
	jobs, err := jobsService.GetAllJobsByUploadID(upload.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(jobs))

	// Mark the job as failed
	err = jobsService.ProcessJobStatus(job.ID, models.JobStatusFailed)
	assert.NoError(t, err)

	// There should now be no one additional forward job associated with the upload
	jobs, err = jobsService.GetAllJobsByUploadID(upload.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(jobs))
}

func TestProcessJobStatusBuildAndAutopkgTestAndForward(t *testing.T) {
	servicesTest := servicestest.SetupTest(t)
	jobsService := servicesTest.Services.Jobs()
	artifactsService := servicesTest.Services.Artifacts()

	// Create an upload
	upload, err := servicesTest.DB.CreateUpload(22, "testsource", "testversion", "testmaintainer", "testchangedby", true, true)
	assert.NoError(t, err)
	assert.NotNil(t, upload)

	// Create a build job for this upload
	job, err := jobsService.CreateJob(
		models.JobTypeBuild, "", models.JobParentTypeUpload, upload.ID,
	)
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
	err = jobsService.ProcessJobStatus(job.ID, models.JobStatusSuccess)
	assert.NoError(t, err)

	// There should now be two new autopkgtest jobs associated with the upload
	jobs, err = jobsService.GetAllJobsByUploadID(upload.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(jobs))

	autopkgTestJob := jobs[1]
	assert.Equal(t, models.JobTypeAutopkgtest, autopkgTestJob.Type)
	assert.Equal(t, upload.ID, autopkgTestJob.ParentID)
	assert.Equal(t, "1", autopkgTestJob.Input, "this job's input should be the build job")

	// Mark the autopkgtest job as completed
	err = jobsService.ProcessJobStatus(autopkgTestJob.ID, models.JobStatusSuccess)
	assert.NoError(t, err)

	// There should now be a forward job associated with the upload
	jobs, err = jobsService.GetAllJobsByUploadID(upload.ID)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(jobs))

	forwardJob := jobs[2]
	assert.Equal(t, models.JobTypeForward, forwardJob.Type)

	// Mark the forward job as completed
	err = jobsService.ProcessJobStatus(forwardJob.ID, models.JobStatusSuccess)
	assert.NoError(t, err)

	// There should be no additional jobs associated with the upload
	jobs, err = jobsService.GetAllJobsByUploadID(upload.ID)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(jobs))
}
