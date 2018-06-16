package webpages_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/routertest"
)

func TestArchiveUpgradesGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	request := httptest.NewRequest(http.MethodGet, "/archive-upgrades", nil)
	response := testRouter.ServeHTTP(request)

	_, err := testRouter.Services.Jobs().CreateArchiveUpgrade(421, 766)
	assert.NoError(t, err)

	response = testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(
		t, response.Body.String(),
		"421",
		"766",
	)
}

func TestNewArchiveUpgradePostHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	user := testRouter.Login()

	archiveUpgrades, err := testRouter.AppCtx.JobsService().GetAllArchiveUpgradesByUserID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(archiveUpgrades))

	form := &url.Values{}
	form.Add("source-suite", "testsourcesuite")
	form.Add("target-suite", "testtargetsuite")
	form.Add("package-count", "42")

	response := testRouter.PostForm("/new-archive-upgrade", form)
	assert.Equal(t, http.StatusSeeOther, response.Result().StatusCode)

	archiveUpgrades, err = testRouter.AppCtx.JobsService().GetAllArchiveUpgradesByUserID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(archiveUpgrades))

	archiveUpgrade := archiveUpgrades[0]
	assert.Equal(t, uint(42), archiveUpgrade.PackageCount)
}
