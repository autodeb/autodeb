package app

import (
	"bufio"
	"bytes"
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

// ProcessUpload receives uploaded files
func (app *App) ProcessUpload(uploadParameters *UploadParameters, content io.Reader) error {
	// Clean the file name, ensure that it contains only a file name and
	// that it isn't something shady like ../../filename.txt
	_, uploadFileName := filepath.Split(uploadParameters.Filename)

	// Check if this is a .changes upload
	isChanges := strings.HasSuffix(uploadFileName, ".changes")

	if isChanges {
		return app.processChangesUpload(uploadFileName, content)
	}
	return app.processFileUpload(uploadFileName, content)
}

func (app *App) processChangesUpload(filename string, content io.Reader) error {
	b, err := ioutil.ReadAll(content)
	if err != nil {
		return err
	}

	changes, err := control.ParseChanges(
		bufio.NewReader(bytes.NewReader(b)),
		"",
	)
	if err != nil {
		return err
	} else if len(changes.ChecksumsSha256) < 1 {
		return fmt.Errorf("changes has no Sha256 checksums")
	}

	//Verify that we have all specified files
	//otherwise, immediately reject the upload
	var pendingFileUploads []*models.PendingFileUpload
	for _, file := range changes.ChecksumsSha256 {
		pendingFileUpload, err := app.dataStore.GetPendingFileUpload(file.Filename, file.Hash, false)
		if err != nil {
			return err
		} else if pendingFileUpload == nil {
			return fmt.Errorf("changes refers to unexisting file %s with hash %s", file.Filename, file.Hash)
		}
		pendingFileUploads = append(pendingFileUploads, pendingFileUpload)
	}

	upload, err := app.dataStore.CreateUpload()
	if err != nil {
		return err
	}

	destDir := filepath.Join(app.UploadsDirectory(), fmt.Sprint(upload.ID))

	//Save the .changes file
	if err := writeDataToDestInFS(
		bytes.NewReader(b),
		destDir,
		filename,
		app.dataFS,
	); err != nil {
		return err
	}

	//Move all files to the upload folder and delete the pendingFileUploads
	for _, pendingFileUpload := range pendingFileUploads {

		sourceDir := filepath.Join(app.UploadedFilesDirectory(), fmt.Sprint(pendingFileUpload.ID))
		source := filepath.Join(sourceDir, pendingFileUpload.Filename)
		dest := filepath.Join(app.UploadsDirectory(), fmt.Sprint(upload.ID), pendingFileUpload.Filename)

		if err := app.dataFS.Rename(source, dest); err != nil {
			//At this point, it is too late to back down. We could have deleted
			//a PendingFileUpload already so we better finish moving what we can
			//and just log the error.
			log.Printf("Couldn't move %s to %s\n", source, dest)
		}

		pendingFileUpload.Completed = true
		if err := app.dataStore.UpdatePendingFileUpload(pendingFileUpload); err != nil {
			// Not stopping, same reason as above.
			log.Printf("Couldn't update pendingFileUpload %d\n", pendingFileUpload.ID)
		}
		app.dataFS.RemoveAll(sourceDir)

	}

	return nil
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
