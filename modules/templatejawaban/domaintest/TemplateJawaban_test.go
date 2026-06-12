package domaintest

import (
	"UnpakSiamida/modules/templatejawaban/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTemplateJawaban_TableName(t *testing.T) {
	tj := domain.TemplateJawaban{}
	assert.Equal(t, "template_pilihanv2", tj.TableName())
}

func TestTemplateJawaban_New(t *testing.T) {
	var nilai uint = 4
	res := domain.NewTemplateJawaban(101, "Sangat Baik", &nilai, 0, "local", "ref-123")
	assert.True(t, res.IsSuccess)
	assert.NotNil(t, res.Value)
	assert.Equal(t, uint(101), res.Value.IdTemplatePertanyaan)
	assert.Equal(t, "Sangat Baik", res.Value.Jawaban)
	assert.Equal(t, &nilai, res.Value.Nilai)

	// Failure case
	resFail := domain.NewTemplateJawaban(101, "Sangat Baik", &nilai, 0, "invalid-owner", "ref-123")
	assert.False(t, resFail.IsSuccess)
	assert.Equal(t, domain.InvalidOwner(), resFail.Error)
}

func TestTemplateJawaban_Update(t *testing.T) {
	var nilai uint = 4
	tj := domain.NewTemplateJawaban(101, "Sangat Baik", &nilai, 0, "local", "ref-123").Value

	// Success case
	var newNilai uint = 5
	res := domain.UpdateTemplateJawaban(tj, tj.UUID, 102, "Luar Biasa", &newNilai, 1, "local", "ref-123")
	assert.True(t, res.IsSuccess)
	assert.Equal(t, uint(102), res.Value.IdTemplatePertanyaan)
	assert.Equal(t, "Luar Biasa", res.Value.Jawaban)
	assert.Equal(t, &newNilai, res.Value.Nilai)
	assert.Equal(t, uint(1), res.Value.IsFreeText)

	// Failure case: Nil object
	resNil := domain.UpdateTemplateJawaban(nil, tj.UUID, 102, "Luar Biasa", &newNilai, 1, "local", "ref-123")
	assert.False(t, resNil.IsSuccess)
	assert.Equal(t, domain.EmptyData(), resNil.Error)

	// Failure case: UUID mismatch
	resMismatch := domain.UpdateTemplateJawaban(tj, uuid.New(), 102, "Luar Biasa", &newNilai, 1, "local", "ref-123")
	assert.False(t, resMismatch.IsSuccess)
	assert.Equal(t, domain.InvalidData(), resMismatch.Error)
}

func TestTemplateJawaban_DeleteAndRestore(t *testing.T) {
	var nilai uint = 4
	tj := domain.NewTemplateJawaban(101, "Sangat Baik", &nilai, 0, "local", "ref-123").Value

	// Delete
	resDel := domain.DeleteTemplateJawaban(tj)
	assert.True(t, resDel.IsSuccess)
	assert.NotNil(t, resDel.Value.DeletedAt)

	// Restore
	resRestore := domain.RestoreTemplateJawaban(tj)
	assert.True(t, resRestore.IsSuccess)
	assert.Nil(t, resRestore.Value.DeletedAt)

	// Nil scenarios
	assert.False(t, domain.DeleteTemplateJawaban(nil).IsSuccess)
	assert.False(t, domain.RestoreTemplateJawaban(nil).IsSuccess)
}
