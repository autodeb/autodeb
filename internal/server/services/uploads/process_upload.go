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
	Autopkgtest   bool
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
		upload, err := service.processChangesUpload(uploadFileName, content, uploadParameters)
		return upload, err
	case ".deb":
		return nil, &uploadError{errors.New("only source uploads are accepted"), true}
	default:
		_, err := service.processFileUpload(uploadFileName, content)
		return nil, err
	}

}

func (service *Service) processChangesUpload(filename string, content io.Reader, uploadParameters *UploadParameters) (*models.Upload, error) {
	contentBytes, err := ioutil.ReadAll(content)
	if err != nil {
		return nil, err
	}

	// Find the signer
	signerID, err := service.pgpService.IdentifySigner(
		bytes.NewReader(contentBytes),
	)
	if err != nil {
		return nil, &uploadError{errors.WithMessage(err, "could not identify the signer"), true}
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

	//Process the .changes file upload
	fileUpload, err := service.processFileUpload(filename, bytes.NewReader(contentBytes))
	if err != nil {
		return nil, err
	}

	// Append the .changes file to the file uploads
	fileUploads = append(fileUploads, fileUpload)

	//Create the upload
	upload, err := service.db.CreateUpload(
		signerID,
		changes.Source,
		changes.Version.String(),
		changes.Maintainer,
		changes.ChangedBy,
		uploadParameters.Autopkgtest,
		uploadParameters.ForwardUpload,
	)
	if err != nil {
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

	//Create a build job. We do this at only after moving files
	//because the job could be picked-up immediately.
	if _, err := service.jobsService.CreateBuildJob(upload.ID); err != nil {
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

	destFolder := filepath.Join(
		service.UploadsDirectory(),
		fmt.Sprint(upload.ID),
	)

	destFile := filepath.Join(
		destFolder,
		fileUpload.Filename,
	)

	if _, err := service.fs.Stat(destFolder); os.IsNotExist(err) {
		if err := service.fs.MkdirAll(destFolder, 0755); err != nil {
			return errors.WithMessage(err, "could not create destination directory")
		}
	}

	if err := service.fs.Rename(source, destFile); err != nil {
		return errors.Errorf("could not move %s to %s", source, destFile)
	}

	fileUpload.Completed = true
	fileUpload.UploadID = upload.ID
	if err := service.db.UpdateFileUpload(fileUpload); err != nil {
		return errors.Errorf("could not mark fileUpload %v as completed", fileUpload.ID)
	}

	service.fs.RemoveAll(sourceDir)

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

func (service *Service) processFileUpload(filename string, content io.Reader) (*models.FileUpload, error) {
	// Save the upload to a temp file on the os's filesystem so that we can
	// calculate the shasum
	tmpfileName, err := writeToTempfile(content)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpfileName)

	shasum, err := sha256.Sum256HexFile(tmpfileName)
	if err != nil {
		return nil, err
	}

	fileUpload, err := service.db.CreateFileUpload(filename, shasum, time.Now())
	if err != nil {
		return nil, err
	}

	destDir := filepath.Join(service.UploadedFilesDirectory(), fmt.Sprint(fileUpload.ID))

	// Open the temporary file
	tmpFile, err := os.Open(tmpfileName)
	if err != nil {
		return nil, err
	}
	defer tmpFile.Close()

	// Write the upload in the fileStorage
	if err := writeDataToDestInFS(
		tmpFile,
		destDir,
		filename,
		service.fs,
	); err != nil {
		return nil, err
	}

	return fileUpload, err
}
