package uploadparametersparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitURLPathOnlyPackage(t *testing.T) {
	filename, splitPath, err := splitURLPath("package.changes")

	assert.Nil(t, err)
	assert.Equal(t, "package.changes", filename)
	assert.Empty(t, splitPath)
	assert.NotNil(t, splitPath)
}

func TestSplitURLPathOnlyPackageSlash(t *testing.T) {
	filename, splitPath, err := splitURLPath("/package.changes")

	assert.Nil(t, err)
	assert.Equal(t, "package.changes", filename)
	assert.Empty(t, splitPath)
	assert.NotNil(t, splitPath)
}

func TestSplitURLPathWithParams(t *testing.T) {
	filename, splitPath, err := splitURLPath("param1/value1/package.changes")

	assert.Nil(t, err)
	assert.Equal(t, "package.changes", filename)
	assert.Equal(t, []string{"param1", "value1"}, splitPath)
}

func TestSplitURLPathNothing(t *testing.T) {
	_, _, err := splitURLPath("")
	assert.EqualError(t, err, "upload parameters should atleast contain the filename")
}
