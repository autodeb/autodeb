package jobrunner

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (jobRunner *JobRunner) execAddBuildToRepository(
	ctx context.Context,
	job *models.Job,
	workingDirectory string,
	artifactsDirectory string,
	logFile io.Writer) error {

	// The job input is the name of the repository
	repositoryName := job.Input

	// Create the repository if it does not exist
	if _, err := jobRunner.apiClient.Aptly().GetOrCreateAndPublishRepository(repositoryName); err != nil {
		return errors.WithMessagef(err, "could get or create repository %s", repositoryName)
	}

	buildJob, err := jobRunner.apiClient.GetJob(job.BuildJobID)
	if err != nil {
		return errors.WithMessagef(err, "could not get build job %d", job.BuildJobID)
	}

	// Retrieve the build job's artifacts
	artifacts, err := jobRunner.apiClient.GetJobArtifacts(buildJob.ID)
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
				repositoryName,
			); err != nil {
				return errors.WithMessagef(err, "could not upload package %s to repository %s", artifact.Filename, job.Input)
			}

		}

	}

	// Publish the repository
	if err := jobRunner.apiClient.Aptly().UpdatePublishedRepositoryDefaults(repositoryName); err != nil {
		return errors.WithMessage(err, "could update aptly repository")
	}
	fmt.Fprintf(logFile, "Updated aptly repository\n")

	return nil
}
