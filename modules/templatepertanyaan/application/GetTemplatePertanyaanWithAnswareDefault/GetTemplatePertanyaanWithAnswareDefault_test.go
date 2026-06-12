package application_test

import (
	"context"
	"testing"

	mockBankSoal "UnpakSiamida/modules/banksoal/application/mock"
	domainBankSoal "UnpakSiamida/modules/banksoal/domain"
	"UnpakSiamida/modules/templatepertanyaan/application/GetTemplatePertanyaanWithAnswareDefault"
	"UnpakSiamida/modules/templatepertanyaan/application/mock"
	"UnpakSiamida/modules/templatepertanyaan/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// === Test By UUID Handler ===

func TestGetTemplatePertanyaanWithAnswareDefaultByUuidQueryHandler_Success(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	expectedTp := &domain.TemplatePertanyaanWithAnswareDefault{
		ID:         1,
		UUID:       uid,
		Pertanyaan: "Test Default Pertanyaan With Answer",
	}

	repo.GetDefaultWithAnswareByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaanWithAnswareDefault, error) {
		assert.Equal(t, uid, id)
		return expectedTp, nil
	}

	handler := &application.GetTemplatePertanyaanWithAnswareDefaultByUuidQueryHandler{
		Repo: repo,
	}

	q := application.GetTemplatePertanyaanWithAnswareDefaultByUuidQuery{
		Uuid: uid.String(),
	}

	res, err := handler.Handle(context.Background(), q)
	assert.NoError(t, err)
	assert.Equal(t, expectedTp, res)
}

func TestGetTemplatePertanyaanWithAnswareDefaultByUuidQueryHandler_InvalidUuid(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}

	handler := &application.GetTemplatePertanyaanWithAnswareDefaultByUuidQueryHandler{
		Repo: repo,
	}

	q := application.GetTemplatePertanyaanWithAnswareDefaultByUuidQuery{
		Uuid: "invalid-uuid",
	}

	res, err := handler.Handle(context.Background(), q)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Nil(t, res)
}

func TestGetTemplatePertanyaanWithAnswareDefaultByUuidQueryHandler_NotFound(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	repo.GetDefaultWithAnswareByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaanWithAnswareDefault, error) {
		return nil, gorm.ErrRecordNotFound
	}

	handler := &application.GetTemplatePertanyaanWithAnswareDefaultByUuidQueryHandler{
		Repo: repo,
	}

	q := application.GetTemplatePertanyaanWithAnswareDefaultByUuidQuery{
		Uuid: uid.String(),
	}

	res, err := handler.Handle(context.Background(), q)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Nil(t, res)
}

// === Test By BankSoal Handler ===

func TestGetTemplatePertanyaanWithAnswareDefaultByBankSoalQueryHandler_Success(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	repoBankSoal := &mockBankSoal.MockRepository{}
	uidBS := uuid.New()

	banksoal := &domainBankSoal.BankSoal{
		ID:    100,
		UUID:  uidBS,
		Judul: "BankSoal Test",
	}

	repoBankSoal.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domainBankSoal.BankSoal, error) {
		assert.Equal(t, uidBS, id)
		return banksoal, nil
	}

	expectedList := []domain.TemplatePertanyaanWithAnswareDefault{
		{
			ID:         1,
			Pertanyaan: "Pertanyaan 1",
		},
	}

	repo.GetDefaultWithAnswareByBankSoalFunc = func(ctx context.Context, id_banksoal uint) ([]domain.TemplatePertanyaanWithAnswareDefault, error) {
		assert.Equal(t, uint(100), id_banksoal)
		return expectedList, nil
	}

	handler := &application.GetTemplatePertanyaanWithAnswareDefaultByBankSoalQueryHandler{
		Repo:         repo,
		RepoBankSoal: repoBankSoal,
	}

	q := application.GetTemplatePertanyaanWithAnswareDefaultByBankSoalQuery{
		UuidBankSoal: uidBS.String(),
	}

	res, err := handler.Handle(context.Background(), q)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.Total)
	assert.Equal(t, expectedList, res.Data)
}

func TestGetTemplatePertanyaanWithAnswareDefaultByBankSoalQueryHandler_InvalidUuid(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	repoBankSoal := &mockBankSoal.MockRepository{}

	handler := &application.GetTemplatePertanyaanWithAnswareDefaultByBankSoalQueryHandler{
		Repo:         repo,
		RepoBankSoal: repoBankSoal,
	}

	q := application.GetTemplatePertanyaanWithAnswareDefaultByBankSoalQuery{
		UuidBankSoal: "invalid-uuid",
	}

	res, err := handler.Handle(context.Background(), q)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kategori is invalid") // Note: domain error returned is InvalidBankSoal()
	assert.Empty(t, res.Data)
}

func TestGetTemplatePertanyaanWithAnswareDefaultByBankSoalQueryHandler_BankSoalNotFound(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	repoBankSoal := &mockBankSoal.MockRepository{}
	uidBS := uuid.New()

	repoBankSoal.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domainBankSoal.BankSoal, error) {
		return nil, gorm.ErrRecordNotFound
	}

	handler := &application.GetTemplatePertanyaanWithAnswareDefaultByBankSoalQueryHandler{
		Repo:         repo,
		RepoBankSoal: repoBankSoal,
	}

	q := application.GetTemplatePertanyaanWithAnswareDefaultByBankSoalQuery{
		UuidBankSoal: uidBS.String(),
	}

	res, err := handler.Handle(context.Background(), q)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kategori is not found") // NotFoundBankSoal returns "kategori is not found"
	assert.Empty(t, res.Data)
}
