package domaintest

import (
	"UnpakSiamida/modules/kuesioner/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKuesioner_TableName(t *testing.T) {
	k := domain.Kuesioner{}
	assert.Equal(t, "kuesionerv2", k.TableName())

	kj := domain.KuesionerJawaban{}
	assert.Equal(t, "kuesioner_jawabanv2", kj.TableName())
}

func TestKuesioner_New(t *testing.T) {
	// Success case
	res := domain.NewKuesioner(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		"soal-1",
		"2024-01-01 10:20:30",
		"local",
		"ref-123",
	)
	assert.True(t, res.IsSuccess)
	assert.NotNil(t, res.Value)
	assert.Equal(t, "soal-1", res.Value.IdBankSoal)

	// Failure case (Invalid Date)
	resFail := domain.NewKuesioner(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		"soal-1",
		"invalid-date",
		"local",
		"ref-123",
	)
	assert.False(t, resFail.IsSuccess)
	assert.Contains(t, resFail.Error.Description, "tanggal mulai")
}

func TestKuesionerJawaban_New(t *testing.T) {
	var idJawaban uint = 500
	freeText := "Free text answer"
	createdBy := "student"
	createdByRef := "npm-123"

	res := domain.NewKuesionerJawaban(
		10,
		20,
		&idJawaban,
		&freeText,
		&createdBy,
		&createdByRef,
	)

	assert.True(t, res.IsSuccess)
	assert.NotNil(t, res.Value)
	assert.Equal(t, uint(10), res.Value.IdKuesioner)
	assert.Equal(t, uint(20), res.Value.IdTemplatePertanyaan)
	assert.Equal(t, &idJawaban, res.Value.IdTemplateJawaban)
	assert.Equal(t, &freeText, res.Value.FreeText)
}
