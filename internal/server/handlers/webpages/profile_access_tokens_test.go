package webpages_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/routertest"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccessTokenPostHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	user := testRouter.Login()

	accessTokens, err := testRouter.AppCtx.TokensService().GetUserTokens(user.ID)
	assert.NotNil(t, accessTokens)
	assert.Equal(t, 0, len(accessTokens))
	assert.NoError(t, err)

	form := &url.Values{}
	form.Add("name", "testname")

	response := testRouter.PostForm("/profile/create-access-token", form)

	assert.Equal(t, http.StatusSeeOther, response.Result().StatusCode)

	accessTokens, err = testRouter.AppCtx.TokensService().GetUserTokens(user.ID)
	assert.NotNil(t, accessTokens)
	assert.Equal(t, 1, len(accessTokens))
	assert.NoError(t, err)

	accessToken := accessTokens[0]
	assert.Equal(t, "testname", accessToken.Name)
	assert.Equal(t, user.ID, accessToken.UserID)
}

func TestRemoveAccessTokenPostHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	user := testRouter.Login()

	token, err := testRouter.AppCtx.TokensService().CreateToken(user.ID, "testname")
	assert.NoError(t, err)

	accessTokens, err := testRouter.AppCtx.TokensService().GetUserTokens(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, accessTokens)
	assert.Equal(t, 1, len(accessTokens))

	form := &url.Values{}
	form.Add("tokenid", fmt.Sprint(token.ID))

	response := testRouter.PostForm("/profile/remove-access-token", form)

	assert.Equal(t, http.StatusSeeOther, response.Result().StatusCode)

	accessTokens, err = testRouter.AppCtx.TokensService().GetUserTokens(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, accessTokens)
	assert.Equal(t, 0, len(accessTokens))
}
