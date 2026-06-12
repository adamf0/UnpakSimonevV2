package infrastructure

import (
	"errors"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type DummyRequest struct {
	Name string
}

func DummyValidator(req DummyRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required.Error("name cannot be blank")),
	)
}

func TestValidationRegistry(t *testing.T) {
	// Register validator
	RegisterValidation(DummyValidator, "Dummy.Error")

	// Get validator
	req := DummyRequest{Name: "Hello"}
	entry, ok := GetValidator(req)
	require.True(t, ok)
	assert.Equal(t, "Dummy.Error", entry.label)

	// Validate: Success case
	err := Validate(req)
	assert.NoError(t, err)

	// Validate: Failure case (normal validation error)
	reqFail := DummyRequest{Name: ""}
	errFail := Validate(reqFail)
	assert.Error(t, errFail)

	// Verify the error is of type ResponseError
	respErr, ok := errFail.(*ResponseError)
	require.True(t, ok)
	assert.Equal(t, "Dummy.Error", respErr.Code)
	
	// Message should contain map with lowercase field name
	msgs, ok := respErr.Message.(map[string]string)
	require.True(t, ok)
	assert.Equal(t, "name cannot be blank", msgs["name"])

	// Test non-ozzo errors wrapping
	RegisterValidation(func(r DummyRequest) error {
		return errors.New("generic error")
	}, "Generic.Error")

	errGeneric := Validate(req)
	assert.Error(t, errGeneric)
	respErrGen, ok := errGeneric.(*ResponseError)
	require.True(t, ok)
	assert.Equal(t, "generic error", respErrGen.Message)
}
