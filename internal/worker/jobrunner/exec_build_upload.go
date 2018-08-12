package jobrunner

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"pault.ag/go/debian/control"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/exec"
	"salsa.debian.org/autodeb-team/autodeb/internal/ftpmasterapi"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (jobRunner *JobRunner) execBuildUpload(
	ctx context.Context,
	job *models.Job,
	workingDirectory string,
	artifactsDirectory string,
	logFile io.Writer) error {

	// Error if we are not building an upload
	if job.ParentType != models.JobParentTypeUpload {
		return errors.Errorf("unsupported parent type %s", job.ParentType)
	}

	// Download all upload artifacts, locating the dsc and debian source
	fileUploads, err := jobRunner.apiClient.GetUploadFiles(job.ParentID)
	if err != nil {
		return errors.WithMessage(err, "could not get upload files")
	}

	var dscPath string
	for _, file := range fileUploads {

		fileContent, err := jobRunner.apiClient.GetUploadFile(file.UploadID, file.Filename)
		if err != nil {
			return errors.WithMessagef(err, "could not get upload file %s from upload id %d", file.Filename, file.UploadID)
		}

		destPath := filepath.Join(workingDirectory, file.Filename)
		if filepath.Ext(destPath) == ".dsc" {
			dscPath = destPath
		}

		destFile, err := os.Create(destPath)
		if err != nil {
			return errors.WithMessagef(err, "could not create dest file at %s", destPath)
		}
		defer destFile.Close()

		if _, err := io.Copy(destFile, fileContent); err != nil {
			return errors.WithMessage(err, "could not copy file")
		}

		destFile.Close()
	}

	// Parse the dsc
	dsc, err := control.ParseDscFile(dscPath)
	if err != nil {
		return errors.WithMessagef(err, "could not parse dsc at %s", dscPath)
	}

	// Ensure that we have all files referred by the .dsc, downloading the missing ones
	// from the archive.
	for _, file := range dsc.ChecksumsSha256 {

		filePath := filepath.Join(workingDirectory, file.Filename)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			continue
		}

		fileContent, err := ftpmasterapi.NewClient(http.DefaultClient).GetFileBySHA256Sum(file.Hash)
		if err != nil {
			return errors.WithMessagef(err, "could not file file with sha256sum %s in the archive", file.Hash)
		}
		defer fileContent.Close()

		destFile, err := os.Create(filePath)
		if err != nil {
			return errors.WithMessagef(err, "could not create %s", filePath)
		}
		defer destFile.Close()

		if _, err := io.Copy(destFile, fileContent); err != nil {
			return errors.WithMessage(err, "could not write file contents")
		}

		destFile.Close()
	}

	// Extract the source package
	if err := exec.RunCtxDirStdoutStderr(
		ctx, workingDirectory, logFile, logFile,
		"dpkg-source", "--extract", dscPath,
	); err != nil {
		return errors.WithMessagef(err, "dpkg-source could not extract %s", dscPath)
	}

	// Find the source directory
	dirs, err := getDirectories(workingDirectory)
	if err != nil {
		return err
	}
	if numDirs := len(dirs); numDirs != 1 {
		return errors.Errorf("multiple directories, cannot guess which one contains the source: %s", dirs)
	}
	sourceDirectory := filepath.Join(workingDirectory, dirs[0])

	// Run sbuild
	if err := exec.RunCtxDirStdoutStderr(
		ctx, sourceDirectory, logFile, logFile,
		"sbuild",
		"--no-clean-source",
		"--nolog",
		"--arch-all",
		"--source",
	); err != nil {
		return errors.WithMessage(err, "sbuild failed")
	}

	// Copy .dsc and referenced files to artifacts directory
	if err := dsc.Copy(artifactsDirectory); err != nil {
		return errors.WithMessage(err, "could not copy dsc and related files to the artifacts directory")
	}

	// Find .changes file
	changes, err := getFirstChangesInDirectory(workingDirectory)
	if err != nil {
		return errors.WithMessage(err, "couldn't get changes file in output directory")
	} else if changes == nil {
		return errors.New("no changes file found in output directory")
	}

	// Copy .changes and referenced files to artifacts directory
	if err := changes.Copy(artifactsDirectory); err != nil {
		return errors.WithMessage(err, "could not copy changes and referenced files to the artifacts directory")
	}

	return nil
}

// getFirstChangesInDirectory returns the first .changes file found in a directory
func getFirstChangesInDirectory(directory string) (*control.Changes, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {

		if filepath.Ext(file.Name()) != ".changes" {
			continue
		}

		changes, err := control.ParseChangesFile(
			filepath.Join(directory, file.Name()),
		)
		return changes, err

	}

	return nil, nil
}

//getDirectories returns a list of all directories in a directory
func getDirectories(dir string) ([]string, error) {
	var directories []string

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.Mode().IsDir() {
			directories = append(directories, file.Name())
		}
	}

	return directories, nil
}
