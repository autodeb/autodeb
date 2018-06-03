package jobrunner

import (
	"context"
	"io"
	"os"
	"path"
	"path/filepath"

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

	// Get the .dsc URL
	dscURL := jobRunner.apiClient.GetUploadDSCURL(job.UploadID)
	dscFileName := path.Base(dscURL.EscapedPath())

	// Download the source
	if err := exec.RunCtxDirStdoutStderr(
		ctx, workingDirectory, logFile, logFile,
		"dget", "--allow-unauthenticated", dscURL.String(),
	); err != nil {
		return errors.WithMessage(err, "dget failed")
	}

	// Get the artifact (deb) that we should test
	artifact, err := jobRunner.apiClient.GetArtifact(job.InputArtifactID)
	if err != nil {
		return errors.WithMessage(err, "could not get job input artifact")
	}

	// Get the deb
	artifactContent, err := jobRunner.apiClient.GetArtifactContent(artifact.ID)
	if err != nil {
		return errors.WithMessage(err, "could not get the job input artifact content")
	}

	// Write the deb to disk
	debPath := filepath.Join(workingDirectory, artifact.Filename)
	deb, err := os.Create(debPath)
	if err != nil {
		return errors.WithMessage(err, "could not create deb")
	}
	defer deb.Close()
	if _, err := io.Copy(deb, artifactContent); err != nil {
		return errors.WithMessage(err, "could not copt artifact content to deb")
	}

	// Run Autopkgtest
	if err := exec.RunCtxDirStdoutStderr(
		ctx, workingDirectory, logFile, logFile,
		"autopkgtest",
		"--no-built-binaries",
		artifact.Filename,
		dscFileName,
		"--",
		"schroot",
		"unstable-amd64-sbuild",
	); err != nil {
		return errors.WithMessage(err, "autopkgtest failed")
	}

	return nil
}
