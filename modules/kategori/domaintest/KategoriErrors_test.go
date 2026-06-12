package domaintest

import (
	"UnpakSiamida/modules/kategori/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKategoriErrors(t *testing.T) {
	assert.Equal(t, "Kategori.EmptyData", domain.EmptyData().Code)
	assert.Equal(t, "Kategori.InvalidUuid", domain.InvalidUuid().Code)
	assert.Equal(t, "Kategori.InvalidData", domain.InvalidData().Code)
	assert.Equal(t, "Kategori.NotFound", domain.NotFound("123").Code)
	assert.Equal(t, "Kategori.InvalidHierarchy", domain.InvalidHierarchy().Code)
	assert.Equal(t, "Kategori.InvalidOwner", domain.InvalidOwner().Code)
}
