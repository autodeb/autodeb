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

	// Get the ArchiveUpgrade
	archiveUpgrade, err := jobRunner.apiClient.GetArchiveUpgrade(job.ParentID)
	if err != nil {
		return errors.WithMessage(err, "could not get the archive job's archive upgrade")
	}

	// Get the jobs
	builds, err := jobRunner.apiClient.GetArchiveUpgradeSuccessfulBuilds(archiveUpgrade.ID)
	if err != nil {
		return errors.WithMessage(err, "could not get archive upgrade jobs")
	}

	// Find the successfull successfull builds
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

				// Add it to the archive rebuild's repository
				if err := jobRunner.apiClient.Aptly().UploadPackageAndAddToRepository(
					artifact.Filename,
					artifactContent,
					archiveUpgrade.RepositoryName(),
				); err != nil {
					return errors.WithMessagef(err, "could not upload package %s to repository %s", artifact.Filename, archiveUpgrade.RepositoryName())
				}

			}

		}

	}

	// Update the repository
	if err := jobRunner.apiClient.Aptly().UpdatePublishedRepositoryDefaults(archiveUpgrade.RepositoryName()); err != nil {
		return errors.WithMessage(err, "could update aptly repository")
	}
	fmt.Fprintf(logFile, "Updated aptly repository\n")

	return nil
}
