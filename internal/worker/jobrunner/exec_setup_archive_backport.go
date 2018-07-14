package jobrunner

import (
	"context"
	"io"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (jobRunner *JobRunner) execSetupArchiveBackport(
	ctx context.Context,
	job *models.Job,
	workingDirectory string,
	artifactsDirectory string,
	logFile io.Writer) error {

	if job.ParentType != models.JobParentTypeArchiveBackport {
		return errors.Errorf("unsupported parent type %s", job.ParentType)
	}

	return nil
}
