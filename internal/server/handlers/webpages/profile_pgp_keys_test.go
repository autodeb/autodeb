package webpages_test

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/pgp"
	"salsa.debian.org/autodeb-team/autodeb/internal/pgp/pgptest"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/routertest"

	"github.com/stretchr/testify/assert"
)

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

func TestRemovePGPKeyPostHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	user := testRouter.Login()
	testRouter.AddPGPKeyToUser(user)

	keys, err := testRouter.AppCtx.PGPService().GetUserPGPKeys(user.ID)
	assert.NotNil(t, keys)
	assert.Equal(t, 1, len(keys))
	assert.NoError(t, err)

	form := &url.Values{}
	form.Add("keyid", fmt.Sprint(keys[0].ID))

	response := testRouter.PostForm("/profile/remove-pgp-key", form)

	assert.Equal(t, http.StatusSeeOther, response.Result().StatusCode)

	keys, err = testRouter.AppCtx.PGPService().GetUserPGPKeys(user.ID)
	assert.NotNil(t, keys)
	assert.Equal(t, 0, len(keys))
	assert.NoError(t, err)
}
