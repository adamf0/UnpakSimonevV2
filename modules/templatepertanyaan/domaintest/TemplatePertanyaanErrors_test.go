package domaintest

import (
	"UnpakSiamida/modules/templatepertanyaan/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplatePertanyaanErrors(t *testing.T) {
	assert.Equal(t, "TemplatePertanyaan.EmptyData", domain.EmptyData().Code)
	assert.Equal(t, "TemplatePertanyaan.InvalidUuid", domain.InvalidUuid().Code)
	assert.Equal(t, "TemplatePertanyaan.InvalidKategori", domain.InvalidKategori().Code)
	assert.Equal(t, "TemplatePertanyaan.InvalidStatus", domain.InvalidStatus().Code)
	assert.Equal(t, "TemplatePertanyaan.InvalidBankSoal", domain.InvalidBankSoal().Code)
	assert.Equal(t, "TemplatePertanyaan.NotFoundKategori", domain.NotFoundKategori().Code)
	assert.Equal(t, "TemplatePertanyaan.NotFoundBankSoal", domain.NotFoundBankSoal().Code)
	assert.Equal(t, "TemplatePertanyaan.InvalidData", domain.InvalidData().Code)
	assert.Equal(t, "TemplatePertanyaan.NotFound", domain.NotFound("123").Code)
	assert.Equal(t, "TemplatePertanyaan.InvalidOwner", domain.InvalidOwner().Code)
}
