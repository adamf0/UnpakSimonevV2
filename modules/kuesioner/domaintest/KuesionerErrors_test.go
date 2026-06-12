package domaintest

import (
	"UnpakSiamida/modules/kuesioner/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKuesionerErrors(t *testing.T) {
	assert.Equal(t, "Kuesioner.EmptyData", domain.EmptyData().Code)
	assert.Equal(t, "Kuesioner.InvalidUuid", domain.InvalidUuid().Code)
	assert.Equal(t, "Kuesioner.InvalidBankSoal", domain.InvalidBankSoal().Code)
	assert.Equal(t, "Kuesioner.InvalidPertanyaan", domain.InvalidPertanyaan().Code)
	assert.Equal(t, "Kuesioner.InvalidJawaban", domain.InvalidJawaban().Code)
	assert.Equal(t, "Kuesioner.NotFound", domain.NotFound("123").Code)
	assert.Equal(t, "Kuesioner.NotFoundBankSoal", domain.NotFoundBankSoal().Code)
	assert.Equal(t, "Kuesioner.NotFoundPertanyaan", domain.NotFoundPertanyaan().Code)
	assert.Equal(t, "Kuesioner.NotFoundJawaban", domain.NotFoundJawaban().Code)
	assert.Equal(t, "Kuesioner.NotFoundResource", domain.NotFoundResource().Code)
	assert.Equal(t, "Kuesioner.NoInfoAccount", domain.NoInfoAccount().Code)
	assert.Equal(t, "Kuesioner.RespondentOnly", domain.RespondentOnly().Code)
	assert.Equal(t, "Kuesioner.InvalidDate", domain.InvalidDate("start").Code)
	assert.Equal(t, "Kuesioner.InvalidDateRange", domain.InvalidDateRange().Code)
	assert.Equal(t, "Kuesioner.Expired", domain.Expired().Code)
}
