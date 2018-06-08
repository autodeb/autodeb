package auth_test

import (
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/routertest"

	"github.com/stretchr/testify/assert"
)

func TestAuthTokenUnauthenticated(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient

	// Create a user and a token
	tokenUser := testRouter.GetOrCreateTestUser()
	token := testRouter.AddTokenToUser(tokenUser)

	// Try to retrieve the current user without the token
	user, err := apiClient.GetCurrentUser()
	assert.NoError(t, err)
	assert.Nil(t, user)

	// Try to retrieve the current user with the token
	apiClient.SetToken(token.Token)
	user, err = apiClient.GetCurrentUser()
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, tokenUser, user)
}
