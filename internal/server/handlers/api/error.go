package api

import (
	"bytes"
	"encoding/json"
)

//Error is the error format that is returned to the user
type Error struct {
	Message string `json:"message"`
}

// JSONError creates a JSON Error object from the input message. If the message
// cannot be encoded, it will be returned as-is so that we still can return
// something to the user.
func JSONError(message string) string {
	e := &Error{
		Message: message,
	}

	b, err := json.Marshal(e)
	if err != nil {
		return message
	}

	return string(b)
}

//ErrorFromJSON returns an Error from json
func ErrorFromJSON(b []byte) (*Error, error) {
	decoder := json.NewDecoder(
		bytes.NewReader(b),
	)
	decoder.DisallowUnknownFields()

	var e Error
	err := decoder.Decode(&e)

	return &e, err
}
