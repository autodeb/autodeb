package aptly

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
)

// UploadFileInDirectory will upload a file in a directory
func (client *APIClient) UploadFileInDirectory(content io.Reader, filename, directory string) error {

	body := &bytes.Buffer{}
	multiPartWriter := multipart.NewWriter(body)
	defer multiPartWriter.Close()

	fileWriter, err := multiPartWriter.CreateFormFile("file", filename)
	if err != nil {
		return errors.WithMessagef(err, "could not create fileWriter for %s", filename)
	}

	if _, err := io.Copy(fileWriter, content); err != nil {
		return errors.WithMessage(err, "could not copy file to filewriter")
	}

	contentType := multiPartWriter.FormDataContentType()

	multiPartWriter.Close()

	resp, err := client.do(
		http.MethodPost,
		fmt.Sprintf("/files/%s", directory),
		contentType,
		body,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	return nil
}
