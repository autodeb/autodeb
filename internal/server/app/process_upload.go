package app

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"salsa.debian.org/autodeb-team/autodeb/internal/crypto/sha256"
)

// UploadParameters defines upload behaviour
type UploadParameters struct {
	Filename      string
	ForwardUpload bool
}

// ProcessUpload receives uploaded files
func (app *App) ProcessUpload(uploadParameters *UploadParameters, uploadContent io.Reader) error {
	// Clean the file name, ensure that it contains only a file name and
	// that it isn't something shady like ../../filename.txt
	_, uploadFileName := filepath.Split(uploadParameters.Filename)

	// Check if this is a .changes upload
	isChanges := strings.HasSuffix(uploadFileName, ".changes")

	// Save the upload to a temp file
	tmpfile, err := writeToTempfile(uploadContent)
	if err != nil {
		return err
	}

	// Find out the destination directory
	var destDir string
	if isChanges {
		upload, err := app.dataStore.CreateUpload()
		if err != nil {
			return err
		}
		destDir = filepath.Join(app.UploadsDirectory(), fmt.Sprint(upload.ID))
	} else {
		shasum, err := sha256.Sum256HexFile(tmpfile)
		if err != nil {
			return err
		}
		destDir = filepath.Join(app.UploadedFilesDirectory(), shasum)
	}

	// move the file to destination
	err = moveFileToDest(tmpfile, destDir, uploadFileName)

	return err
}

func writeToTempfile(data io.Reader) (string, error) {
	// Create temp file
	tmpfile, err := ioutil.TempFile("", "autodeb-upload")
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()

	// Get the file name
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

func moveFileToDest(source, destdir, destfilename string) error {
	fulldestpath := filepath.Join(destdir, destfilename)

	// Create the destination directory
	if err := os.MkdirAll(destdir, 0744); err != nil {
		return err
	}

	// Move the file
	err := os.Rename(source, fulldestpath)

	return err
}
