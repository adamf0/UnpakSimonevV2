package domaintest

import (
	"UnpakSiamida/modules/prodi/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProdiErrors(t *testing.T) {
	errEmpty := domain.EmptyData()
	assert.Equal(t, "Prodi.EmptyData", errEmpty.Code)

	errNotFound := domain.NotFound("IPA")
	assert.Equal(t, "Prodi.NotFound", errNotFound.Code)
	assert.Contains(t, errNotFound.Description, "IPA")
}
