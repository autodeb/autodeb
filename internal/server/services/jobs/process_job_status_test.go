package jobs_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services/servicestest"
)

func TestProcessJobStatusUpgradeAutopkgtest(t *testing.T) {
	servicesTest := servicestest.SetupTest(t)
	jobsService := servicesTest.Services.Jobs()

	archiveUpgrade, err := jobsService.CreateArchiveUpgrade(1, 33)
	assert.NoError(t, err)

	// Create a package upgrade job in the context of an archive upgrade
	upgradeJob, err := jobsService.CreateJob(
		models.JobTypePackageUpgrade, "", models.JobParentTypeArchiveUpgrade, archiveUpgrade.ID,
	)
	assert.NoError(t, err)
	assert.NotNil(t, upgradeJob)

	// There should be two jobs associated with the archive upgrade:
	//  - the upgrade init
	//  - the package upgrade job
	jobs, err := jobsService.GetAllJobsByArchiveUpgradeID(archiveUpgrade.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(jobs))

	// Mark the job as successfull
	err = jobsService.ProcessJobStatus(upgradeJob.ID, models.JobStatusSuccess)
	assert.NoError(t, err)

	// There should now be a new autopkgtest job associated with the archive upgrade
	jobs, err = jobsService.GetAllJobsByArchiveUpgradeID(archiveUpgrade.ID)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(jobs))

	autopkgTestJob := jobs[2]
	assert.Equal(t, models.JobTypeAutopkgtest, autopkgTestJob.Type)
	assert.Equal(t, models.JobParentTypeArchiveUpgrade, autopkgTestJob.ParentType)
	assert.Equal(t, archiveUpgrade.ID, autopkgTestJob.ParentID)
	assert.Equal(t, fmt.Sprint(upgradeJob.ID), autopkgTestJob.Input, "this job's input should be the package upgrade job")

	// Mark the autopkgtest job as completed
	err = jobsService.ProcessJobStatus(autopkgTestJob.ID, models.JobStatusSuccess)
	assert.NoError(t, err)

	// There should be no new jobs created
	jobs, err = jobsService.GetAllJobsByArchiveUpgradeID(archiveUpgrade.ID)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(jobs))
}

func TestProcessJobStatusUploadBuildAndDontForward(t *testing.T) {
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

func TestProcessJobStatusUploadBuildAndAutopkgTestAndForward(t *testing.T) {
	servicesTest := servicestest.SetupTest(t)
	jobsService := servicesTest.Services.Jobs()

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

	// Mark the job as successfull
	err = jobsService.ProcessJobStatus(job.ID, models.JobStatusSuccess)
	assert.NoError(t, err)

	// There should now be a new autopkgtest job associated with the upload
	jobs, err = jobsService.GetAllJobsByUploadID(upload.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(jobs))

	autopkgTestJob := jobs[1]
	assert.Equal(t, models.JobTypeAutopkgtest, autopkgTestJob.Type)
	assert.Equal(t, upload.ID, autopkgTestJob.ParentID)
	assert.Equal(t, fmt.Sprint(job.ID), autopkgTestJob.Input, "this job's input should be the build job")

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
