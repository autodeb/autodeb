package uploadparametersparser

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNoParameters(t *testing.T) {
	r := httptest.NewRequest("PUT", "/package.changes", strings.NewReader("test"))

	parameters, err := Parse(r)

	assert.NotNil(t, parameters)
	assert.Nil(t, err)
}

func TestParseUnrecognizedPathParameters(t *testing.T) {
	r := httptest.NewRequest("PUT", "/aa/true/package.changes", strings.NewReader("test"))

	parameters, err := Parse(r)

	assert.Nil(t, parameters)
	assert.EqualError(t, err, "unrecognized upload parameter: aa")
}

func TestParseUnrecognizedQueryParameters(t *testing.T) {
	r := httptest.NewRequest("PUT", "/package.changes?aa=bb", strings.NewReader("test"))

	parameters, err := Parse(r)

	assert.Nil(t, parameters)
	assert.EqualError(t, err, "unrecognized upload parameter: aa")
}
