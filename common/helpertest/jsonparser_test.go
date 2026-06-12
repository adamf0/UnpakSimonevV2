package commontest

import (
	"UnpakSiamida/common/helper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRejectDuplicateKeys(t *testing.T) {
	// Valid JSON without duplicate keys
	validJSON := []byte(`{"name": "john", "age": 30, "address": {"city": "New York", "zip": "10001"}}`)
	assert.NoError(t, helper.RejectDuplicateKeys(validJSON))

	// Invalid JSON: Duplicate keys at top level
	dupJSON := []byte(`{"name": "john", "age": 30, "name": "doe"}`)
	err := helper.RejectDuplicateKeys(dupJSON)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate field detected")

	// Invalid JSON: Duplicate keys in nested object
	nestedDupJSON := []byte(`{"name": "john", "info": {"age": 30, "age": 31}}`)
	err2 := helper.RejectDuplicateKeys(nestedDupJSON)
	assert.Error(t, err2)
	assert.Contains(t, err2.Error(), "duplicate field detected")

	// Invalid JSON: Weird characters in field name
	badFieldJSON := []byte(`{"na-me": "john"}`)
	err3 := helper.RejectDuplicateKeys(badFieldJSON)
	assert.Error(t, err3)
	assert.Contains(t, err3.Error(), "invalid field name")
}
