package jobrunner

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"time"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/udd"
)

func (jobRunner *JobRunner) execSetupArchiveUpgrade(
	ctx context.Context,
	job *models.Job,
	workingDirectory string,
	artifactsDirectory string,
	logFile io.Writer) error {

	if job.ParentType != models.JobParentTypeArchiveUpgrade {
		return errors.Errorf("unsupported parent type %s", job.ParentType)
	}

	// Get the archive upgrade
	archiveUpgrade, err := jobRunner.apiClient.GetArchiveUpgrade(job.ParentID)
	if err != nil {
		return errors.WithMessage(err, "could not get archive upgrade")
	}

	// Create aptly repository
	repo, err := jobRunner.apiClient.Aptly().GetOrCreateAndPublishRepository(
		archiveUpgrade.RepositoryName(),
	)
	if err != nil {
		return errors.WithMessage(err, "could not create aptly repository")
	}
	fmt.Fprintf(logFile, "Created repository %+v\n", repo)

	// Get all packages that need upgrading
	sourcePackages, err := udd.PackagesWithNewerUpstreamVersions()
	if err != nil {
		return errors.WithMessagef(err, "could not get source packages to update")
	} else if len(sourcePackages) < 1 {
		return errors.New("there are no packages to upgrade in the archive")
	}

	rand.Seed(time.Now().UnixNano())

	fmt.Fprintln(logFile, "Creating upgrade jobs...")

	for i := 0; i < archiveUpgrade.PackageCount || archiveUpgrade.PackageCount < 0; i++ {

		if len(sourcePackages) < 1 {
			fmt.Fprintln(logFile, "there are no more source pacakges to ugprade...")
			break
		}

		// select an index
		index := rand.Intn(len(sourcePackages))

		// grab it from the package list
		pkg := sourcePackages[index]

		// remove it from the package list
		sourcePackages = append(sourcePackages[:index], sourcePackages[index+1:]...)

		// Create an upgrade job
		fmt.Fprintf(
			logFile,
			"Creating job: %s %s => %s... ",
			pkg.Package,
			pkg.DebianUversion,
			pkg.UpstreamVersion,
		)

		job, err := jobRunner.apiClient.CreateJob(
			&models.Job{
				Type:       models.JobTypePackageUpgrade,
				Input:      pkg.Package,
				ParentType: models.JobParentTypeArchiveUpgrade,
				ParentID:   archiveUpgrade.ID,
			},
		)
		if err != nil {
			return errors.WithMessagef(err, "could not create upgrade job for package %s", pkg.Package)
		}

		fmt.Fprintf(logFile, "job id #%d\n", job.ID)
	}

	return nil
}
