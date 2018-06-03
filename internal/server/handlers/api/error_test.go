package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/api"
)

func TestJSONError(t *testing.T) {
	jsonError := api.JSONError("hey")

	assert.Equal(t, `{"message":"hey"}`, jsonError)
}

func TestErrorFromJSON(t *testing.T) {
	e, err := api.ErrorFromJSON(
		[]byte(`{"message":"hey"}`),
	)

	assert.NoError(t, err)
	assert.Equal(t, "hey", e.Message)
}
