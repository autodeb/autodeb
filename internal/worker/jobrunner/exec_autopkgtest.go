package jobrunner

import (
	"context"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/exec"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (jobRunner *JobRunner) execAutopkgtest(
	ctx context.Context,
	job *models.Job,
	workingDirectory string,
	artifactsDirectory string,
	logFile io.Writer) error {

	if job.ParentType != models.JobParentTypeUpload {
		return errors.Errorf("unsupported parent type %s", job.ParentType)
	}

	// The job input is a build job
	buildJobID, err := strconv.Atoi(job.Input)
	if err != nil {
		return errors.WithMessage(err, "could not convert input to int")
	}

	// Get the .dsc URL
	dscURL := jobRunner.apiClient.GetUploadDSCURL(job.ParentID)
	dscFileName := path.Base(dscURL.EscapedPath())

	// Download the source
	if err := exec.RunCtxDirStdoutStderr(
		ctx, workingDirectory, logFile, logFile,
		"dget", "--allow-unauthenticated", dscURL.String(),
	); err != nil {
		return errors.WithMessage(err, "dget failed")
	}

	// Get the artifacts of the build job
	artifacts, err := jobRunner.apiClient.GetJobArtifacts(uint(buildJobID))
	if err != nil {
		return errors.WithMessage(err, "could not get job artifacts")
	}

	var debFilenames []string

	// Get the artifacts (debs) that we should test
	for _, artifact := range artifacts {

		if filepath.Ext(artifact.Filename) != ".deb" {
			continue
		}

		// Get the deb
		artifactContent, err := jobRunner.apiClient.GetArtifactContent(artifact.ID)
		if err != nil {
			return errors.WithMessage(err, "could not get the artifact content")
		}

		// Write it on disk
		debPath := filepath.Join(workingDirectory, artifact.Filename)
		deb, err := os.Create(debPath)
		if err != nil {
			return errors.WithMessage(err, "could not create deb")
		}
		defer deb.Close()
		if _, err := io.Copy(deb, artifactContent); err != nil {
			return errors.WithMessage(err, "could not copt artifact content to deb")
		}

		debFilenames = append(debFilenames, artifact.Filename)
	}

	args := []string{
		"--no-built-binaries",
	}
	args = append(
		args,
		debFilenames...,
	)
	args = append(
		args,
		dscFileName,
		"--",
		"schroot",
		"unstable-amd64-sbuild",
	)

	// Run Autopkgtest
	if err := exec.RunCtxDirStdoutStderr(
		ctx, workingDirectory, logFile, logFile,
		"autopkgtest", args...,
	); err != nil {
		// TODO: Autopkgtest's exit code 8 means that no tests were run.
		// We shouldn't error when that happens.
		return errors.WithMessage(err, "autopkgtest failed")
	}

	return nil
}
