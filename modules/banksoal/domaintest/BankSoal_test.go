package domaintest

import (
	"UnpakSiamida/common/helper"
	"UnpakSiamida/modules/banksoal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBankSoal_TableName(t *testing.T) {
	bs := domain.BankSoal{}
	assert.Equal(t, "bank_soalv2", bs.TableName())

	bse := domain.BankSoalExt{}
	assert.Equal(t, "bank_soal_extendv2", bse.TableName())
}

func TestBankSoal_New(t *testing.T) {
	// Success case
	res := domain.NewBankSoal("Judul Soal", helper.StrPtr("content"), helper.StrPtr("deskripsi"), helper.StrPtr("20241"), "local", "ref-123")
	assert.True(t, res.IsSuccess)
	assert.NotNil(t, res.Value)
	assert.Equal(t, "Judul Soal", res.Value.Judul)
	assert.Equal(t, "draf", res.Value.Status)

	// Failure case: invalid owner
	resFail := domain.NewBankSoal("Judul Soal", nil, nil, nil, "invalid-owner", "ref-123")
	assert.False(t, resFail.IsSuccess)
	assert.Equal(t, domain.InvalidOwner(), resFail.Error)
}

func TestBankSoal_Update(t *testing.T) {
	bs := domain.NewBankSoal("Judul Soal", nil, nil, nil, "local", "ref-123").Value

	// Success case
	res := domain.UpdateBankSoal(bs, bs.UUID, "Updated Title", helper.StrPtr("new content"), nil, nil, "local", "ref-123")
	assert.True(t, res.IsSuccess)
	assert.Equal(t, "Updated Title", res.Value.Judul)

	// Failure case: Nil object
	resNil := domain.UpdateBankSoal(nil, bs.UUID, "Updated Title", nil, nil, nil, "local", "ref-123")
	assert.False(t, resNil.IsSuccess)
	assert.Equal(t, domain.EmptyData(), resNil.Error)

	// Failure case: UUID mismatch
	resMismatch := domain.UpdateBankSoal(bs, uuid.New(), "Updated Title", nil, nil, nil, "local", "ref-123")
	assert.False(t, resMismatch.IsSuccess)
	assert.Equal(t, domain.InvalidData(), resMismatch.Error)

	// Failure case: Invalid owner
	resOwner := domain.UpdateBankSoal(bs, bs.UUID, "Updated Title", nil, nil, nil, "invalid-owner", "ref-123")
	assert.False(t, resOwner.IsSuccess)
	assert.Equal(t, domain.InvalidOwner(), resOwner.Error)
}

func TestBankSoal_UpdateTime(t *testing.T) {
	bs := domain.NewBankSoal("Judul Soal", nil, nil, nil, "local", "ref-123").Value

	// Success case
	start := "2024-01-01"
	end := "2024-01-10"
	res := domain.UpdateTimeBankSoal(bs, bs.UUID, &start, &end)
	assert.True(t, res.IsSuccess)
	assert.NotNil(t, res.Value.TanggalMulai)
	assert.NotNil(t, res.Value.TanggalAkhir)

	// Failure case: Nil object
	resNil := domain.UpdateTimeBankSoal(nil, bs.UUID, &start, &end)
	assert.False(t, resNil.IsSuccess)
	assert.Equal(t, domain.EmptyData(), resNil.Error)

	// Failure case: UUID mismatch
	resMismatch := domain.UpdateTimeBankSoal(bs, uuid.New(), &start, &end)
	assert.False(t, resMismatch.IsSuccess)
	assert.Equal(t, domain.InvalidData(), resMismatch.Error)

	// Failure case: Invalid start date format
	badStart := "invalid-date"
	resBadStart := domain.UpdateTimeBankSoal(bs, bs.UUID, &badStart, &end)
	assert.False(t, resBadStart.IsSuccess)
	assert.Equal(t, domain.InvalidDate("tanggal awal"), resBadStart.Error)

	// Failure case: Invalid end date format
	badEnd := "invalid-date"
	resBadEnd := domain.UpdateTimeBankSoal(bs, bs.UUID, &start, &badEnd)
	assert.False(t, resBadEnd.IsSuccess)
	assert.Equal(t, domain.InvalidDate("tanggal akhir"), resBadEnd.Error)

	// Failure case: Overlapping/invalid range
	resOverlap := domain.UpdateTimeBankSoal(bs, bs.UUID, &start, &badStart) // fails on date format
	assert.False(t, resOverlap.IsSuccess)

	overlapEnd := "2023-12-31"
	resOverlap2 := domain.UpdateTimeBankSoal(bs, bs.UUID, &start, &overlapEnd)
	assert.False(t, resOverlap2.IsSuccess)
	assert.Equal(t, domain.InvalidDateRange(), resOverlap2.Error)
}

func TestBankSoal_DeleteRestoreResetTime(t *testing.T) {
	bs := domain.NewBankSoal("Judul Soal", nil, nil, nil, "local", "ref-123").Value

	// Delete
	resDel := domain.DeleteBankSoal(bs)
	assert.True(t, resDel.IsSuccess)
	assert.NotNil(t, resDel.Value.DeletedAt)

	// Restore
	resRestore := domain.RestoreBankSoal(bs)
	assert.True(t, resRestore.IsSuccess)
	assert.Nil(t, resRestore.Value.DeletedAt)

	// Reset Time
	start := "2024-01-01"
	domain.UpdateTimeBankSoal(bs, bs.UUID, &start, nil)
	assert.NotNil(t, bs.TanggalMulai)

	resReset := domain.ResetTimeBankSoal(bs)
	assert.True(t, resReset.IsSuccess)
	assert.Nil(t, resReset.Value.TanggalMulai)

	// Nil scenarios
	assert.False(t, domain.DeleteBankSoal(nil).IsSuccess)
	assert.False(t, domain.RestoreBankSoal(nil).IsSuccess)
	assert.False(t, domain.ResetTimeBankSoal(nil).IsSuccess)
}

func TestBankSoal_ChangeStatus(t *testing.T) {
	bs := domain.NewBankSoal("Judul Soal", nil, nil, nil, "local", "ref-123").Value

	// Success status
	res := domain.ChangeStatus(bs, "active")
	assert.True(t, res.IsSuccess)
	assert.Equal(t, "active", res.Value.Status)

	// Invalid status
	resFail := domain.ChangeStatus(bs, "invalid-status")
	assert.False(t, resFail.IsSuccess)
	assert.Equal(t, domain.InvalidStatus(), resFail.Error)

	// Nil scenario
	assert.False(t, domain.ChangeStatus(nil, "active").IsSuccess)
}

func TestBankSoal_Copy(t *testing.T) {
	bs := domain.NewBankSoal("Original", nil, nil, nil, "local", "ref-123").Value

	// Copy count = 0
	resCopy1 := domain.CopyBankSoal(bs, 0, "local", "ref-123")
	assert.True(t, resCopy1.IsSuccess)
	assert.Equal(t, "salin - Original", resCopy1.Value.Judul)

	// Copy count = 1
	resCopy2 := domain.CopyBankSoal(bs, 1, "local", "ref-123")
	assert.True(t, resCopy2.IsSuccess)
	assert.Equal(t, "salin (2) - Original", resCopy2.Value.Judul)

	// Nil scenario
	assert.False(t, domain.CopyBankSoal(nil, 0, "local", "").IsSuccess)
}

func TestBankSoalExt_AddTimeBankSoalExt(t *testing.T) {
	bs := domain.NewBankSoal("Judul Soal", nil, nil, nil, "local", "ref-123").Value
	bs.ID = 12

	start := "2024-05-01"
	end := "2024-05-15"

	t.Run("Success", func(t *testing.T) {
		res := domain.AddTimeBankSoalExt(bs, bs.UUID, &start, &end, "local", "ref-123")
		assert.True(t, res.IsSuccess)
		assert.NotNil(t, res.Value)
		assert.Equal(t, uint(12), res.Value.IdBankSoal)
		assert.Equal(t, "local", *res.Value.CreatedBy)
	})

	t.Run("Fail - nil banksoal", func(t *testing.T) {
		res := domain.AddTimeBankSoalExt(nil, bs.UUID, &start, &end, "local", "ref-123")
		assert.False(t, res.IsSuccess)
		assert.Equal(t, domain.InvalidData(), res.Error)
	})

	t.Run("Fail - invalid owner", func(t *testing.T) {
		res := domain.AddTimeBankSoalExt(bs, bs.UUID, &start, &end, "invalid-owner", "ref-123")
		assert.False(t, res.IsSuccess)
		assert.Equal(t, domain.InvalidOwner(), res.Error)
	})

	t.Run("Fail - bad start date format", func(t *testing.T) {
		bad := "invalid-date"
		res := domain.AddTimeBankSoalExt(bs, bs.UUID, &bad, &end, "local", "ref-123")
		assert.False(t, res.IsSuccess)
		assert.Equal(t, domain.InvalidDate("tanggal awal"), res.Error)
	})

	t.Run("Fail - bad end date format", func(t *testing.T) {
		bad := "invalid-date"
		res := domain.AddTimeBankSoalExt(bs, bs.UUID, &start, &bad, "local", "ref-123")
		assert.False(t, res.IsSuccess)
		assert.Equal(t, domain.InvalidDate("tanggal akhir"), res.Error)
	})

	t.Run("Fail - date overlap", func(t *testing.T) {
		overlapEnd := "2024-04-30"
		res := domain.AddTimeBankSoalExt(bs, bs.UUID, &start, &overlapEnd, "local", "ref-123")
		assert.False(t, res.IsSuccess)
		assert.Equal(t, domain.InvalidDateRange(), res.Error)
	})
}
