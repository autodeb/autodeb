package jobrunner

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (jobRunner *JobRunner) execCreateArchiveUpgradeRepository(
	ctx context.Context,
	job *models.Job,
	workingDirectory string,
	artifactsDirectory string,
	logFile io.Writer) error {

	if job.ParentType != models.JobParentTypeArchiveUpgrade {
		return errors.Errorf("unsupported parent type %s", job.ParentType)
	}

	// Get the jobs
	builds, err := jobRunner.apiClient.GetArchiveUpgradeSuccessfulBuilds(job.ParentID)
	if err != nil {
		return errors.WithMessage(err, "could not get archive upgrade jobs")
	}

	// Find the successfull successfull builds
	// TODO: create an API endpoint for this so that we don't have to duplicate this logic all over.
	for _, build := range builds {

		// Retrieve the build job's artifacts
		artifacts, err := jobRunner.apiClient.GetJobArtifacts(build.ID)
		if err != nil {
			return errors.WithMessage(err, "could not retrieve the job's artifacts")
		}

		// Find debs and submit them to the repository
		for _, artifact := range artifacts {

			if filepath.Ext(artifact.Filename) == ".deb" {

				fmt.Fprintf(logFile, "Adding %s to the repository...\n", artifact.Filename)

				// Get the .deb
				artifactContent, err := jobRunner.apiClient.GetArtifactContent(artifact.ID)
				if err != nil {
					return errors.WithMessagef(err, "could not retrieve the content for artifact id %d", artifact.ID)
				}

				uploadDir := fmt.Sprintf("archive-upgrade-%d", job.ParentID)

				// Upload it to aptly
				if err := jobRunner.apiClient.Aptly().UploadFileInDirectory(
					artifactContent,
					artifact.Filename,
					uploadDir,
				); err != nil {
					return errors.WithMessagef(err, "could not upload %s to aptly", artifact.Filename)
				}

				// Add the package
				if err := jobRunner.apiClient.Aptly().AddPackageToRepository(
					artifact.Filename,
					uploadDir,
					fmt.Sprintf("archive-upgrade-%d", job.ParentID),
				); err != nil {
					return errors.WithMessagef(err, "could not add %s to the repository", artifact.Filename)
				}

			}

		}

	}

	return nil
}
