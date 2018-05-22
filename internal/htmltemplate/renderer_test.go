package htmltemplate

import (
	"testing"

	"github.com/stretchr/testify/require"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
)

func setupTest() (*Renderer, filesystem.FS) {
	fs := filesystem.NewMemMapFS()
	renderer := NewRenderer(fs, true)
	return renderer, fs
}

func TestTemplatesAreCached(t *testing.T) {
	renderer, fs := setupTest()
	fs.Create("parent")
	fs.Create("child")

	createdTemplate, err := renderer.getOrCreateTemplate("parent", "child")
	require.Nil(t, err)
	require.NotNil(t, createdTemplate)

	// not using require.Equal because we want to compare pointers!
	cachedTemplate, err := renderer.getOrCreateTemplate("parent", "child")
	require.Nil(t, err)
	require.NotNil(t, cachedTemplate)
	require.True(t, createdTemplate == cachedTemplate, "the previously created template should be returned")
}

func TestTemplateNames(t *testing.T) {
	renderer, fs := setupTest()
	fs.Create("parent")
	fs.Create("child")

	_, ok := renderer.cache.m["parent+child"]
	require.False(t, ok, "the key should not exist yet")
	require.Equal(t, 0, len(renderer.cache.m), "the cache should be empty")

	_, err := renderer.getOrCreateTemplate("parent", "child")
	require.Nil(t, err)

	_, ok = renderer.cache.m["parent+child"]
	require.True(t, ok, "the key should have been created")
	require.Equal(t, 1, len(renderer.cache.m), "only one template should have been created")
}
