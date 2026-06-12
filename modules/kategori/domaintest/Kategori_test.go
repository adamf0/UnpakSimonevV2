package domaintest

import (
	"UnpakSiamida/modules/kategori/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestKategori_TableName(t *testing.T) {
	k := domain.Kategori{}
	assert.Equal(t, "kategoriv2", k.TableName())
}

func TestKategori_New(t *testing.T) {
	// Success case
	res := domain.NewKategori("Kategori Test", nil, "local", "ref-123")
	assert.True(t, res.IsSuccess)
	assert.NotNil(t, res.Value)
	assert.Equal(t, "Kategori Test", res.Value.NamaKategori)
	assert.Nil(t, res.Value.SubKategori)
	assert.Equal(t, "local", *res.Value.CreatedBy)

	// Failure case (Invalid Owner)
	resFail := domain.NewKategori("Kategori Test", nil, "invalid-owner", "ref-123")
	assert.False(t, resFail.IsSuccess)
	assert.Equal(t, domain.InvalidOwner(), resFail.Error)
}

func TestKategori_Update(t *testing.T) {
	kategori := domain.NewKategori("Kategori Test", nil, "local", "ref-123").Value

	// Success case
	res := domain.UpdateKategori(kategori, kategori.UUID, "Updated Name", nil, "local", "ref-456")
	assert.True(t, res.IsSuccess)
	assert.Equal(t, "Updated Name", res.Value.NamaKategori)

	// Failure case: Nil object
	resNil := domain.UpdateKategori(nil, kategori.UUID, "Updated Name", nil, "local", "ref-456")
	assert.False(t, resNil.IsSuccess)
	assert.Equal(t, domain.EmptyData(), resNil.Error)

	// Failure case: UUID mismatch
	resMismatch := domain.UpdateKategori(kategori, uuid.New(), "Updated Name", nil, "local", "ref-456")
	assert.False(t, resMismatch.IsSuccess)
	assert.Equal(t, domain.InvalidData(), resMismatch.Error)
}

func TestKategori_Move(t *testing.T) {
	kategori := domain.NewKategori("Kategori Test", nil, "local", "ref-123").Value
	kategori.ID = 10

	// Success case
	var parentID uint = 20
	res := domain.MoveKategori(kategori, &parentID)
	assert.True(t, res.IsSuccess)
	assert.Equal(t, &parentID, res.Value.SubKategori)

	// Failure case: Same ID hierarchy loop
	resSelf := domain.MoveKategori(kategori, &kategori.ID)
	assert.False(t, resSelf.IsSuccess)
	assert.Equal(t, domain.InvalidHierarchy(), resSelf.Error)

	// Failure case: Nil object
	resNil := domain.MoveKategori(nil, &parentID)
	assert.False(t, resNil.IsSuccess)
	assert.Equal(t, domain.EmptyData(), resNil.Error)
}

func TestKategori_DeleteAndRestore(t *testing.T) {
	kategori := domain.NewKategori("Kategori Test", nil, "local", "ref-123").Value

	// Delete
	resDel := domain.DeleteKategori(kategori)
	assert.True(t, resDel.IsSuccess)
	assert.NotNil(t, resDel.Value.DeletedAt)

	// Restore
	resRestore := domain.RestoreKategori(kategori)
	assert.True(t, resRestore.IsSuccess)
	assert.Nil(t, resRestore.Value.DeletedAt)

	// Nil scenarios
	assert.False(t, domain.DeleteKategori(nil).IsSuccess)
	assert.False(t, domain.RestoreKategori(nil).IsSuccess)
}

func TestKategori_Copy(t *testing.T) {
	kategori := domain.NewKategori("Original", nil, "local", "ref-123").Value

	// Copy count = 0
	resCopy1 := domain.CopyKategori(kategori, 0, "local", "ref-123")
	assert.True(t, resCopy1.IsSuccess)
	assert.Equal(t, "salin - Original", resCopy1.Value.NamaKategori)

	// Copy count = 1
	resCopy2 := domain.CopyKategori(kategori, 1, "local", "ref-123")
	assert.True(t, resCopy2.IsSuccess)
	assert.Equal(t, "salin (2) - Original", resCopy2.Value.NamaKategori)

	// Nil scenario
	assert.False(t, domain.CopyKategori(nil, 0, "local", "ref-123").IsSuccess)
}
