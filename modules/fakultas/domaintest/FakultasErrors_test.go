package domaintest

import (
	"UnpakSiamida/modules/fakultas/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakultasErrors(t *testing.T) {
	errEmpty := domain.EmptyData()
	assert.Equal(t, "Fakultas.EmptyData", errEmpty.Code)

	errNotFound := domain.NotFound("FKIP")
	assert.Equal(t, "Fakultas.NotFound", errNotFound.Code)
	assert.Contains(t, errNotFound.Description, "FKIP")
}
