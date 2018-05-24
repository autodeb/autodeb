package uploadqueue_test

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/endpoints/api"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/routertest"
)

func TestProcessFileUpload(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	testAppCtx := testRouter.AppCtx
	fs := testAppCtx.UploadsService().FS()
	db := testRouter.DB

	_, err := fs.Stat(filepath.Join(testAppCtx.UploadsService().UploadedFilesDirectory(), "1"))
	require.Error(t, err, "the file directory should not exist")

	request, _ := http.NewRequest(
		http.MethodPut,
		"/upload/test.dsc",
		strings.NewReader("this is a test file\n"),
	)

	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusCreated, response.Result().StatusCode)
	assert.Equal(t, "", response.Body.String())

	_, err = fs.Stat(filepath.Join(testAppCtx.UploadsService().UploadedFilesDirectory(), "1"))
	assert.NoError(t, err)

	_, err = fs.Stat(filepath.Join(testAppCtx.UploadsService().UploadedFilesDirectory(), "1", "test.dsc"))
	assert.NoError(t, err)

	expectedSHASum := "b6668cf8c46c7075e18215d922e7812ca082fa6cc34668d00a6c20aee4551fb6"

	fileUpload, err := db.GetFileUploadByFileNameSHASumCompleted(
		"test.dsc",
		expectedSHASum,
		false,
	)
	assert.NoError(t, err)
	assert.NotNil(t, fileUpload)
	assert.Equal(t, uint(0), fileUpload.UploadID)

	assert.Equal(t, uint(1), fileUpload.ID)
	assert.Equal(t, "test.dsc", fileUpload.Filename)
	assert.Equal(t, expectedSHASum, fileUpload.SHA256Sum)
	assert.Equal(t, false, fileUpload.Completed)
}

func TestUploadDebRejected(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	request, _ := http.NewRequest(
		http.MethodPut,
		"/upload/test.deb",
		strings.NewReader("this is a deb\n"),
	)

	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	apiErr, err := api.ErrorFromJSON(response.Body.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, "only source uploads are accepted", apiErr.Message)
}

const dummyChangesFile = `-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA512

Format: 1.8
Date: Wed, 04 Apr 2018 14:28:29 -0400
Source: autodeb
Binary: autodeb-server autodeb-worker
Architecture: source
Version: 1.0+ds1-1
Distribution: unstable
Urgency: medium
Maintainer: Alexandre Viau <aviau@debian.org>
Changed-By: Changed By <changed.by@debian.org>
Description:
 autodeb-server - main autodeb server
 autodeb-worker - autodeb worker component
Changes:
 autodeb (1.0+ds1-1) unstable; urgency=medium
 .
   * Less bugs.
Checksums-Sha1:
 804d716fc5844f1cc5516c8f0be7a480517fdea2 20 test.dsc
Checksums-Sha256:
 b6668cf8c46c7075e18215d922e7812ca082fa6cc34668d00a6c20aee4551fb6 20 test.dsc
Files:
 66ad00916013ea0f7a6550f762b1de1d 20 utils optional test.dsc
-----BEGIN PGP SIGNATURE-----

iQEzBAEBCgAdFiEEi18odTLPQt8c9/fq0x67v8LkDPMFAlsHDq4ACgkQ0x67v8Lk
DPPunQf/fTGl2oB0idOmbt//Xem/Gbn/OX7IAHgxBpqe4p/j/zALUr5uSxj7wkS0
MsbUEgwerzLkceT9yp5ous1T0hiqd1FLy9YvIBGilxPAXm5iNICCHbdC8xX0zrm8
yXHG700/5Z6RsQU0YhokktvUhgxcdRqg9ujII6OgoEVqiYt9a0UynpaGTrCSfd5Z
MQ1vG5UTx5P1H/O107bTRZIHeRdXy45xytyZvLQRkbLw1A7/iWNiP2sySYYFNQRu
/ynKZTDDvPTEBWnAUyCNOu9hCvXhAhSHPDp6fqx1Qvp6jg78sqSAiWKitLBfWW5k
pUA+9T5iTl+RDUR356uU4G+n8mWQ+w==
=8onE
-----END PGP SIGNATURE-----
`

func TestProcessChangesBadFormatRejected(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	request, _ := http.NewRequest(
		http.MethodPut,
		"/upload/test.changes",
		strings.NewReader("test"),
	)

	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	_, err := api.ErrorFromJSON(response.Body.Bytes())
	assert.NoError(t, err)
}

func TestProcessChangesMissingFile(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	user := testRouter.GetOrCreateTestUser()
	testRouter.AddPGPKeyToUser(user)

	request, _ := http.NewRequest(
		http.MethodPut,
		"/upload/test.changes",
		strings.NewReader(dummyChangesFile),
	)

	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	apiErr, err := api.ErrorFromJSON(response.Body.Bytes())
	assert.NoError(t, err)
	assert.Contains(t, apiErr.Message, "changes refers to unexisting file test.dsc")
}

func TestProcessChangesUnsigned(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	request, _ := http.NewRequest(
		http.MethodPut,
		"/upload/test.changes",
		strings.NewReader("unsigned stuff"),
	)

	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)

	apiErr, err := api.ErrorFromJSON(response.Body.Bytes())
	assert.NoError(t, err)
	assert.Contains(t, apiErr.Message, "could not identify the signer")
}

func TestProcessChanges(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	testAppCtx := testRouter.AppCtx
	fs := testAppCtx.UploadsService().FS()
	db := testRouter.DB

	user := testRouter.GetOrCreateTestUser()
	testRouter.AddPGPKeyToUser(user)

	request, _ := http.NewRequest(
		http.MethodPut,
		"/upload/test.dsc",
		strings.NewReader("this is a test file\n"),
	)

	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusCreated, response.Result().StatusCode)
	assert.Equal(t, "", response.Body.String())

	request, _ = http.NewRequest(
		http.MethodPut,
		"/upload/test.changes",
		strings.NewReader(dummyChangesFile),
	)

	response = testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusCreated, response.Result().StatusCode)
	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))

	var upload models.Upload
	json.Unmarshal(response.Body.Bytes(), &upload)

	assert.Equal(t, uint(1), upload.ID)
	assert.Equal(t, "autodeb", upload.Source)
	assert.Equal(t, "1.0+ds1-1", upload.Version)
	assert.Equal(t, "Alexandre Viau <aviau@debian.org>", upload.Maintainer)
	assert.Equal(t, "Changed By <changed.by@debian.org>", upload.ChangedBy)

	_, err := fs.Stat(filepath.Join(testAppCtx.UploadsService().UploadedFilesDirectory(), "1"))
	assert.Error(t, err, "the uploaded files directory should be removed")

	_, err = fs.Stat(filepath.Join(testAppCtx.UploadsService().UploadsDirectory(), "1"))
	assert.NoError(t, err)

	_, err = fs.Stat(filepath.Join(testAppCtx.UploadsService().UploadsDirectory(), "1", "test.changes"))
	assert.NoError(t, err)

	_, err = fs.Stat(filepath.Join(testAppCtx.UploadsService().UploadsDirectory(), "1", "test.dsc"))
	assert.NoError(t, err)

	fileUpload, err := db.GetFileUpload(uint(1))
	assert.NoError(t, err)
	assert.Equal(t, true, fileUpload.Completed)
	assert.Equal(t, uint(1), fileUpload.UploadID)

	jobs, err := db.GetAllJobs()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(jobs))

	job := jobs[0]
	assert.Equal(t, uint(1), job.ID)
	assert.Equal(t, uint(1), job.UploadID)
	assert.Equal(t, models.JobTypeBuild, job.Type)
	assert.Equal(t, models.JobStatusQueued, job.Status)
}
