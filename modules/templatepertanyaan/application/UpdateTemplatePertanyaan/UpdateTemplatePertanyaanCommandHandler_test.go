package application_test

import (
	"context"
	"testing"

	mockBankSoal "UnpakSiamida/modules/banksoal/application/mock"
	domainBankSoal "UnpakSiamida/modules/banksoal/domain"
	mockKategori "UnpakSiamida/modules/kategori/application/mock"
	domainKategori "UnpakSiamida/modules/kategori/domain"
	"UnpakSiamida/modules/templatepertanyaan/application/UpdateTemplatePertanyaan"
	"UnpakSiamida/modules/templatepertanyaan/application/mock"
	"UnpakSiamida/modules/templatepertanyaan/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUpdateTemplatePertanyaanCommandHandler_Success(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	repoKategori := &mockKategori.MockKategoriRepository{}
	repoBankSoal := &mockBankSoal.MockRepository{}

	uuidTP := uuid.New()
	uuidKategori := uuid.New()
	uuidBankSoal := uuid.New()

	existingTP := &domain.TemplatePertanyaan{
		ID:           1,
		UUID:         uuidTP,
		IdBankSoal:   2,
		Pertanyaan:   "Pertanyaan Lama",
		JenisPilihan: "pilihan_ganda",
		Bobot:        1,
		Required:     0,
	}

	kategori := &domainKategori.Kategori{
		ID:           10,
		UUID:         uuidKategori,
		NamaKategori: "Kategori Baru",
	}

	banksoal := &domainBankSoal.BankSoal{
		ID:    20,
		UUID:  uuidBankSoal,
		Judul: "BankSoal Baru",
	}

	repo.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaan, error) {
		assert.Equal(t, uuidTP, id)
		return existingTP, nil
	}

	repoKategori.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domainKategori.Kategori, error) {
		assert.Equal(t, uuidKategori, id)
		return kategori, nil
	}

	repoBankSoal.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domainBankSoal.BankSoal, error) {
		assert.Equal(t, uuidBankSoal, id)
		return banksoal, nil
	}

	repo.UpdateFunc = func(ctx context.Context, tp *domain.TemplatePertanyaan) error {
		assert.Equal(t, uint(20), tp.IdBankSoal)
		assert.Equal(t, "Pertanyaan Baru", tp.Pertanyaan)
		assert.Equal(t, "essay", tp.JenisPilihan)
		assert.Equal(t, uint(5), tp.Bobot)
		assert.Equal(t, uint(10), *tp.IdKategori)
		assert.Equal(t, 1, tp.Required)
		return nil
	}

	handler := &application.UpdateTemplatePertanyaanCommandHandler{
		Repo:         repo,
		RepoKategori: repoKategori,
		RepoBankSoal: repoBankSoal,
	}

	cmd := application.UpdateTemplatePertanyaanCommand{
		Uuid:         uuidTP.String(),
		UuidBankSoal: uuidBankSoal.String(),
		Pertanyaan:   "Pertanyaan Baru",
		JenisPilihan: "essay",
		Bobot:        "5",
		UuidKategori: uuidKategori.String(),
		Required:     1,
		SID:          "sid-new",
		Resource:     "local",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, uuidTP.String(), res)
}

func TestUpdateTemplatePertanyaanCommandHandler_InvalidUuid(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}

	handler := &application.UpdateTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.UpdateTemplatePertanyaanCommand{
		Uuid:         "invalid-uuid",
		UuidBankSoal: uuid.New().String(),
		UuidKategori: uuid.New().String(),
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "uuid is invalid")
	assert.Empty(t, res)
}

func TestUpdateTemplatePertanyaanCommandHandler_NotFound(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	repoKategori := &mockKategori.MockKategoriRepository{}
	repoBankSoal := &mockBankSoal.MockRepository{}

	uuidTP := uuid.New()
	uuidKategori := uuid.New()
	uuidBankSoal := uuid.New()

	repo.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaan, error) {
		return nil, gorm.ErrRecordNotFound
	}

	handler := &application.UpdateTemplatePertanyaanCommandHandler{
		Repo:         repo,
		RepoKategori: repoKategori,
		RepoBankSoal: repoBankSoal,
	}

	cmd := application.UpdateTemplatePertanyaanCommand{
		Uuid:         uuidTP.String(),
		UuidBankSoal: uuidBankSoal.String(),
		UuidKategori: uuidKategori.String(),
		Bobot:        "5",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Empty(t, res)
}
