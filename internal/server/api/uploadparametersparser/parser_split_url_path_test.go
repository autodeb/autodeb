package uploadparametersparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitURLPathWithParams(t *testing.T) {
	filename, splitPath, err := splitURLPath("param1/value1/package.changes")

	assert.Nil(t, err)
	assert.Equal(t, "package.changes", filename)
	assert.Equal(t, []string{"param1", "value1"}, splitPath)
}
