package app

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"salsa.debian.org/autodeb-team/autodeb/internal/crypto/sha256"
	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"

	"pault.ag/go/debian/control"
)

// UploadParameters defines upload behaviour
type UploadParameters struct {
	Filename      string
	ForwardUpload bool
}

type uploadError struct {
	error
	isInputError bool
}

func (err *uploadError) IsInputError() bool {
	return err.isInputError
}

// ProcessUpload receives uploaded files
func (app *App) ProcessUpload(uploadParameters *UploadParameters, content io.Reader) (*models.Upload, error) {
	// Clean the file name, ensure that it contains only a file name and
	// that it isn't something shady like ../../filename.txt
	_, uploadFileName := filepath.Split(uploadParameters.Filename)

	// Check if this is a .changes upload
	isChanges := strings.HasSuffix(uploadFileName, ".changes")

	if isChanges {
		upload, err := app.processChangesUpload(uploadFileName, content)
		return upload, err
	}

	err := app.processFileUpload(uploadFileName, content)
	return nil, err
}

func (app *App) processChangesUpload(filename string, content io.Reader) (*models.Upload, error) {
	contentBytes, err := ioutil.ReadAll(content)
	if err != nil {
		return nil, err
	}

	//Parse the changes file
	changes, err := control.ParseChanges(
		bufio.NewReader(bytes.NewReader(contentBytes)),
		"",
	)
	if err != nil {
		return nil, &uploadError{err, true}
	}
	if len(changes.ChecksumsSha256) < 1 {
		return nil, &uploadError{errors.New("changes has no SHA256 checksums"), true}
	}

	//Verify that we have all specified files
	//otherwise, immediately reject the upload
	pendingFileUploads, err := app.getChangesPendingFileUploads(changes)
	if err != nil {
		return nil, err
	}

	//Create the upload
	upload, err := app.dataStore.CreateUpload(
		changes.Source,
		changes.Version.String(),
		changes.Maintainer,
		changes.ChangedBy,
	)
	if err != nil {
		return nil, err
	}

	//Save the .changes file
	if err := writeDataToDestInFS(
		bytes.NewReader(contentBytes),
		filepath.Join(app.UploadsDirectory(), fmt.Sprint(upload.ID)),
		filename,
		app.dataFS,
	); err != nil {
		return nil, err
	}

	//Move all files to the upload folder
	for _, pendingFileUpload := range pendingFileUploads {
		if err := app.movePendingFileUpload(pendingFileUpload, upload); err != nil {
			//At this point, it is too late to back down. We could have moved
			//a PendingFileUpload already so we better finish moving what we can
			//and just log the error.
			log.Printf("cannot move file upload: %v\n", err)
		}
	}

	//Create jobs. We do this at only after moving files
	//because the job could be picked-up immediately.
	if _, err := app.dataStore.CreateJob(models.JobTypeBuild, upload.ID); err != nil {
		return nil, err
	}

	return upload, nil
}

//movePendingFileUpload will move a pendingFileUpload to the upload directory
//and mark the pending file upload as completed
func (app *App) movePendingFileUpload(pendingFileUpload *models.PendingFileUpload, upload *models.Upload) error {
	sourceDir := filepath.Join(
		app.UploadedFilesDirectory(),
		fmt.Sprint(pendingFileUpload.ID),
	)

	source := filepath.Join(
		sourceDir,
		pendingFileUpload.Filename,
	)

	dest := filepath.Join(
		app.UploadsDirectory(),
		fmt.Sprint(upload.ID),
		pendingFileUpload.Filename,
	)

	if err := app.dataFS.Rename(source, dest); err != nil {
		return fmt.Errorf("could not move %s to %s", source, dest)
	}

	pendingFileUpload.Completed = true
	if err := app.dataStore.UpdatePendingFileUpload(pendingFileUpload); err != nil {
		return fmt.Errorf("could not mark pendingFileUpload %v as completed", pendingFileUpload.ID)
	}

	app.dataFS.RemoveAll(sourceDir)

	return nil
}

//getChangedPendingFileUploads returns the pending file uploads associated with this changes file
func (app *App) getChangesPendingFileUploads(changes *control.Changes) ([]*models.PendingFileUpload, error) {
	var pendingFileUploads []*models.PendingFileUpload

	for _, file := range changes.ChecksumsSha256 {
		pendingFileUpload, err := app.dataStore.GetPendingFileUpload(
			file.Filename,
			file.Hash,
			false,
		)
		if err != nil {
			return nil, err
		}

		if pendingFileUpload == nil {
			return nil, &uploadError{
				fmt.Errorf(
					"changes refers to unexisting file %s with hash %s",
					file.Filename,
					file.Hash,
				),
				true,
			}
		}

		pendingFileUploads = append(pendingFileUploads, pendingFileUpload)
	}

	return pendingFileUploads, nil
}

func (app *App) processFileUpload(filename string, content io.Reader) error {
	// Save the upload to a temp file on the os's filesystem so that we can
	// calculate the shasum
	tmpfileName, err := writeToTempfile(content)
	if err != nil {
		return err
	}
	defer os.Remove(tmpfileName)

	shasum, err := sha256.Sum256HexFile(tmpfileName)
	if err != nil {
		return err
	}

	pendingFileUpload, err := app.dataStore.CreatePendingFileUpload(filename, shasum, time.Now())
	if err != nil {
		return err
	}

	destDir := filepath.Join(app.UploadedFilesDirectory(), fmt.Sprint(pendingFileUpload.ID))

	// Open the temporary file
	tmpFile, err := os.Open(tmpfileName)
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	// Write the upload in the fileStorage
	err = writeDataToDestInFS(
		tmpFile,
		destDir,
		filename,
		app.dataFS,
	)

	return err
}

func writeDataToDestInFS(data io.Reader, destDir, destFileName string, fs filesystem.FS) error {
	// Create the destination directory
	if err := fs.MkdirAll(destDir, 0744); err != nil {
		return err
	}

	// Create the destination file
	destFile, err := fs.Create(filepath.Join(destDir, destFileName))
	if err != nil {
		fs.RemoveAll(destDir)
		return err
	}
	defer destFile.Close()

	// Write data
	if _, err := io.Copy(destFile, data); err != nil {
		destFile.Close()
		fs.RemoveAll(destDir)
	}

	return nil
}

func writeToTempfile(data io.Reader) (string, error) {
	// Create temp file
	tmpfile, err := ioutil.TempFile("", "autodeb-upload")
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()

	filename := tmpfile.Name()

	// Write data
	_, err = io.Copy(tmpfile, data)
	if err != nil {
		tmpfile.Close()
		os.Remove(filename)
		return "", err
	}

	return filename, nil
}
