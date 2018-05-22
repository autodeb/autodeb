package webpages_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/pgp"
	"salsa.debian.org/autodeb-team/autodeb/internal/pgp/pgptest"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/routertest"

	"github.com/stretchr/testify/assert"
)

func TestProfileGetHandlerAuthenticated(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	user := testRouter.Login()

	request := httptest.NewRequest(http.MethodGet, "/profile", nil)
	response := testRouter.ServeHTTP(request)

	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(t, response.Body.String(), user.Username)

	testRouter.Logout()

	request = httptest.NewRequest(http.MethodGet, "/profile", nil)
	response = testRouter.ServeHTTP(request)

	assert.Equal(t, http.StatusSeeOther, response.Result().StatusCode)
	assert.NotContains(t, response.Body.String(), user.Username)
}

func TestAddPGPKeyPostHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	user := testRouter.Login()

	keys, err := testRouter.AppCtx.PGPService().GetUserPGPKeys(user.ID)
	assert.NotNil(t, keys)
	assert.Equal(t, 0, len(keys))
	assert.NoError(t, err)

	proof, err := pgp.Clearsign(
		strings.NewReader(testRouter.AppCtx.PGPService().ExpectedPGPKeyProofText(user.ID)),
		strings.NewReader(pgptest.TestKeyPrivate),
	)
	assert.NoError(t, err)

	form := &url.Values{}
	form.Add("key", pgptest.TestKeyPublic)
	form.Add("proof", proof)

	response := testRouter.PostForm("/profile/add-pgp-key", form)

	assert.Equal(t, http.StatusSeeOther, response.Result().StatusCode)

	keys, err = testRouter.AppCtx.PGPService().GetUserPGPKeys(user.ID)
	assert.NotNil(t, keys)
	assert.Equal(t, 1, len(keys))
	assert.NoError(t, err)

	key := keys[0]
	assert.Equal(t, pgptest.TestKeyFingerprint, key.Fingerprint)
	assert.Equal(t, pgptest.TestKeyPublic, key.PublicKey)
}
