package domaintest

import (
	"UnpakSiamida/modules/templatepertanyaan/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTemplatePertanyaan_TableName(t *testing.T) {
	tp := domain.TemplatePertanyaan{}
	assert.Equal(t, "template_pertanyaanv2", tp.TableName())
}

func TestTemplatePertanyaan_New(t *testing.T) {
	// Success case
	var kategori uint = 5
	res := domain.NewTemplatePertanyaan(1, "Pertanyaan A?", "radio", 2, &kategori, 1, "local", "ref-123")
	assert.True(t, res.IsSuccess)
	assert.NotNil(t, res.Value)
	assert.Equal(t, "Pertanyaan A?", res.Value.Pertanyaan)
	assert.Equal(t, uint(2), res.Value.Bobot)
	assert.Equal(t, "draf", res.Value.Status)

	// Bobot <= 0 defaults to 1
	resZeroBobot := domain.NewTemplatePertanyaan(1, "Pertanyaan A?", "radio", 0, &kategori, 1, "local", "ref-123")
	assert.True(t, resZeroBobot.IsSuccess)
	assert.Equal(t, uint(1), resZeroBobot.Value.Bobot)

	// Failure case
	resFail := domain.NewTemplatePertanyaan(1, "Pertanyaan A?", "radio", 2, nil, 1, "invalid-owner", "ref-123")
	assert.False(t, resFail.IsSuccess)
	assert.Equal(t, domain.InvalidOwner(), resFail.Error)
}

func TestTemplatePertanyaan_Update(t *testing.T) {
	tp := domain.NewTemplatePertanyaan(1, "Pertanyaan A?", "radio", 2, nil, 1, "local", "ref-123").Value

	// Success case
	res := domain.UpdateTemplatePertanyaan(tp, tp.UUID, 2, "Updated Question", "checkbox", 3, nil, 0, "local", "ref-123")
	assert.True(t, res.IsSuccess)
	assert.Equal(t, "Updated Question", res.Value.Pertanyaan)
	assert.Equal(t, "checkbox", res.Value.JenisPilihan)
	assert.Equal(t, uint(3), res.Value.Bobot)

	// Failure case: Nil object
	resNil := domain.UpdateTemplatePertanyaan(nil, tp.UUID, 2, "Updated Question", "checkbox", 3, nil, 0, "local", "ref-123")
	assert.False(t, resNil.IsSuccess)
	assert.Equal(t, domain.EmptyData(), resNil.Error)

	// Failure case: UUID mismatch
	resMismatch := domain.UpdateTemplatePertanyaan(tp, uuid.New(), 2, "Updated Question", "checkbox", 3, nil, 0, "local", "ref-123")
	assert.False(t, resMismatch.IsSuccess)
	assert.Equal(t, domain.InvalidData(), resMismatch.Error)
}

func TestTemplatePertanyaan_DeleteAndRestore(t *testing.T) {
	tp := domain.NewTemplatePertanyaan(1, "Pertanyaan A?", "radio", 2, nil, 1, "local", "ref-123").Value

	// Delete
	resDel := domain.DeleteTemplatePertanyaan(tp)
	assert.True(t, resDel.IsSuccess)
	assert.NotNil(t, resDel.Value.DeletedAt)

	// Restore
	resRestore := domain.RestoreTemplatePertanyaan(tp)
	assert.True(t, resRestore.IsSuccess)
	assert.Nil(t, resRestore.Value.DeletedAt)

	// Nil scenarios
	assert.False(t, domain.DeleteTemplatePertanyaan(nil).IsSuccess)
	assert.False(t, domain.RestoreTemplatePertanyaan(nil).IsSuccess)
}

func TestTemplatePertanyaan_ChangeStatus(t *testing.T) {
	tp := domain.NewTemplatePertanyaan(1, "Pertanyaan A?", "radio", 2, nil, 1, "local", "ref-123").Value

	// Success status
	res := domain.ChangeStatus(tp, "active")
	assert.True(t, res.IsSuccess)
	assert.Equal(t, "active", res.Value.Status)

	// Invalid status
	resFail := domain.ChangeStatus(tp, "invalid-status")
	assert.False(t, resFail.IsSuccess)
	assert.Equal(t, domain.InvalidStatus(), resFail.Error)

	// Nil scenario
	assert.False(t, domain.ChangeStatus(nil, "active").IsSuccess)
}

func TestTemplatePertanyaan_Copy(t *testing.T) {
	tp := domain.NewTemplatePertanyaan(1, "Original Question", "radio", 2, nil, 1, "local", "ref-123").Value

	// Copy count = 0
	resCopy1 := domain.CopyTemplatePertanyaan(tp, 0, "local", "ref-123")
	assert.True(t, resCopy1.IsSuccess)
	assert.Equal(t, "salin - Original Question", resCopy1.Value.Pertanyaan)

	// Copy count = 1
	resCopy2 := domain.CopyTemplatePertanyaan(tp, 1, "local", "ref-123")
	assert.True(t, resCopy2.IsSuccess)
	assert.Equal(t, "salin (2) - Original Question", resCopy2.Value.Pertanyaan)

	// Nil scenario
	assert.False(t, domain.CopyTemplatePertanyaan(nil, 0, "local", "ref-123").IsSuccess)
}
