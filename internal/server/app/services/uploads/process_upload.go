package uploads

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"salsa.debian.org/autodeb-team/autodeb/internal/crypto/sha256"
	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
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
func (service *Service) ProcessUpload(uploadParameters *UploadParameters, content io.Reader) (*models.Upload, error) {
	// Clean the file name, ensure that it contains only a file name and
	// that it isn't something shady like ../../filename.txt
	_, uploadFileName := filepath.Split(uploadParameters.Filename)

	switch ext := filepath.Ext(uploadFileName); ext {
	case ".changes":
		upload, err := service.processChangesUpload(uploadFileName, content)
		return upload, err
	case ".deb":
		return nil, &uploadError{errors.New("only source uploads are accepted"), true}
	default:
		err := service.processFileUpload(uploadFileName, content)
		return nil, err
	}

}

func (service *Service) processChangesUpload(filename string, content io.Reader) (*models.Upload, error) {
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

	// Verify that the upload refers to at least one file
	if len(changes.ChecksumsSha256) < 1 {
		return nil, &uploadError{errors.New("changes has no SHA256 checksums"), true}
	}

	//Verify that we have all specified files
	//otherwise, immediately reject the upload
	fileUploads, err := service.getChangesFileUploads(changes)
	if err != nil {
		return nil, err
	}

	//Create the upload
	upload, err := service.db.CreateUpload(
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
		filepath.Join(service.UploadsDirectory(), fmt.Sprint(upload.ID)),
		filename,
		service.dataFS,
	); err != nil {
		return nil, err
	}

	//Move all files to the upload folder
	for _, fileUpload := range fileUploads {
		if err := service.moveFileUpload(fileUpload, upload); err != nil {
			//At this point, it is too late to back down. We could have moved
			//a FileUpload already so we better finish moving what we can
			//and just log the error.
			log.Printf("cannot move file upload: %v\n", err)
		}
	}

	//Create jobs. We do this at only after moving files
	//because the job could be picked-up immediately.
	if _, err := service.db.CreateJob(models.JobTypeBuild, upload.ID); err != nil {
		return nil, err
	}

	return upload, nil
}

//moveFileUpload will move a FileUpload to the upload directory
//and mark the FileUpload as completed
func (service *Service) moveFileUpload(fileUpload *models.FileUpload, upload *models.Upload) error {
	sourceDir := filepath.Join(
		service.UploadedFilesDirectory(),
		fmt.Sprint(fileUpload.ID),
	)

	source := filepath.Join(
		sourceDir,
		fileUpload.Filename,
	)

	dest := filepath.Join(
		service.UploadsDirectory(),
		fmt.Sprint(upload.ID),
		fileUpload.Filename,
	)

	if err := service.dataFS.Rename(source, dest); err != nil {
		return errors.Errorf("could not move %s to %s", source, dest)
	}

	fileUpload.Completed = true
	fileUpload.UploadID = upload.ID
	if err := service.db.UpdateFileUpload(fileUpload); err != nil {
		return errors.Errorf("could not mark fileUpload %v as completed", fileUpload.ID)
	}

	service.dataFS.RemoveAll(sourceDir)

	return nil
}

//getChangedFileUploads returns the FileUploads associated with this changes file
func (service *Service) getChangesFileUploads(changes *control.Changes) ([]*models.FileUpload, error) {
	var fileUploads []*models.FileUpload

	for _, file := range changes.ChecksumsSha256 {
		fileUpload, err := service.db.GetFileUploadByFileNameSHASumCompleted(
			file.Filename,
			file.Hash,
			false,
		)
		if err != nil {
			return nil, err
		}

		if fileUpload == nil {
			return nil, &uploadError{
				errors.Errorf(
					"changes refers to unexisting file %s with hash %s",
					file.Filename,
					file.Hash,
				),
				true,
			}
		}

		fileUploads = append(fileUploads, fileUpload)
	}

	return fileUploads, nil
}

func (service *Service) processFileUpload(filename string, content io.Reader) error {
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

	fileUpload, err := service.db.CreateFileUpload(filename, shasum, time.Now())
	if err != nil {
		return err
	}

	destDir := filepath.Join(service.UploadedFilesDirectory(), fmt.Sprint(fileUpload.ID))

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
		service.dataFS,
	)

	return err
}
