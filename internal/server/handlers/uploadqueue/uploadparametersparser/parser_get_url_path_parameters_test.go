package uploadparametersparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUrlPathParameters(t *testing.T) {
	params := getURLPathParameters(
		[]string{
			"param1", "value1",
			"param2", "value2",
			"param1", "value2",
			"paramnovalue",
		},
	)

	expected := map[string][]string{
		"param1": []string{"value1", "value2"},
		"param2": []string{"value2"},
	}

	assert.NotNil(t, params)
	assert.Equal(t, expected, params)
}
