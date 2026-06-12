package application_test

import (
	"context"
	"errors"
	"testing"

	"UnpakSiamida/modules/templatepertanyaan/application/CreateTemplatePertanyaan"
	"UnpakSiamida/modules/templatepertanyaan/application/mock"
	"UnpakSiamida/modules/templatepertanyaan/domain"

	mockBankSoal "UnpakSiamida/modules/banksoal/application/mock"
	domainBankSoal "UnpakSiamida/modules/banksoal/domain"
	mockKategori "UnpakSiamida/modules/kategori/application/mock"
	domainKategori "UnpakSiamida/modules/kategori/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTemplatePertanyaanCommandHandler_Success(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	repoKategori := &mockKategori.MockKategoriRepository{}
	repoBankSoal := &mockBankSoal.MockRepository{}

	uuidKategori := uuid.New()
	uuidBankSoal := uuid.New()

	kategori := &domainKategori.Kategori{
		ID:           10,
		UUID:         uuidKategori,
		NamaKategori: "Kategori Test",
	}

	banksoal := &domainBankSoal.BankSoal{
		ID:    20,
		UUID:  uuidBankSoal,
		Judul: "BankSoal Test",
	}

	repoKategori.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domainKategori.Kategori, error) {
		assert.Equal(t, uuidKategori, id)
		return kategori, nil
	}

	repoBankSoal.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domainBankSoal.BankSoal, error) {
		assert.Equal(t, uuidBankSoal, id)
		return banksoal, nil
	}

	repo.CreateFunc = func(ctx context.Context, tp *domain.TemplatePertanyaan) error {
		assert.Equal(t, uint(20), tp.IdBankSoal)
		assert.Equal(t, "Apakah ini pertanyaan?", tp.Pertanyaan)
		assert.Equal(t, "pilihan_ganda", tp.JenisPilihan)
		assert.Equal(t, uint(5), tp.Bobot)
		assert.Equal(t, uint(10), *tp.IdKategori)
		assert.Equal(t, 1, tp.Required)
		return nil
	}

	handler := &application.CreateTemplatePertanyaanCommandHandler{
		Repo:         repo,
		RepoKategori: repoKategori,
		RepoBankSoal: repoBankSoal,
	}

	cmd := application.CreateTemplatePertanyaanCommand{
		UuidBankSoal: uuidBankSoal.String(),
		Pertanyaan:   "Apakah ini pertanyaan?",
		JenisPilihan: "pilihan_ganda",
		Bobot:        "5",
		UuidKategori: uuidKategori.String(),
		Required:     1,
		SID:          "sid-abc",
		Resource:     "local",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestCreateTemplatePertanyaanCommandHandler_InvalidKategoriUuid(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	repoKategori := &mockKategori.MockKategoriRepository{}
	repoBankSoal := &mockBankSoal.MockRepository{}

	handler := &application.CreateTemplatePertanyaanCommandHandler{
		Repo:         repo,
		RepoKategori: repoKategori,
		RepoBankSoal: repoBankSoal,
	}

	cmd := application.CreateTemplatePertanyaanCommand{
		UuidBankSoal: uuid.New().String(),
		Pertanyaan:   "Apakah ini pertanyaan?",
		JenisPilihan: "pilihan_ganda",
		Bobot:        "5",
		UuidKategori: "invalid-uuid",
		Required:     1,
		SID:          "sid-abc",
		Resource:     "local",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kategori is invalid")
	assert.Empty(t, res)
}

func TestCreateTemplatePertanyaanCommandHandler_InvalidBankSoalUuid(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	repoKategori := &mockKategori.MockKategoriRepository{}
	repoBankSoal := &mockBankSoal.MockRepository{}

	handler := &application.CreateTemplatePertanyaanCommandHandler{
		Repo:         repo,
		RepoKategori: repoKategori,
		RepoBankSoal: repoBankSoal,
	}

	cmd := application.CreateTemplatePertanyaanCommand{
		UuidBankSoal: "invalid-uuid",
		Pertanyaan:   "Apakah ini pertanyaan?",
		JenisPilihan: "pilihan_ganda",
		Bobot:        "5",
		UuidKategori: uuid.New().String(),
		Required:     1,
		SID:          "sid-abc",
		Resource:     "local",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kategori is invalid")
	assert.Empty(t, res)
}

func TestCreateTemplatePertanyaanCommandHandler_KategoriNotFound(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	repoKategori := &mockKategori.MockKategoriRepository{}
	repoBankSoal := &mockBankSoal.MockRepository{}

	uuidKategori := uuid.New()
	uuidBankSoal := uuid.New()

	repoKategori.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domainKategori.Kategori, error) {
		return nil, errors.New("not found")
	}

	handler := &application.CreateTemplatePertanyaanCommandHandler{
		Repo:         repo,
		RepoKategori: repoKategori,
		RepoBankSoal: repoBankSoal,
	}

	cmd := application.CreateTemplatePertanyaanCommand{
		UuidBankSoal: uuidBankSoal.String(),
		Pertanyaan:   "Apakah ini pertanyaan?",
		JenisPilihan: "pilihan_ganda",
		Bobot:        "5",
		UuidKategori: uuidKategori.String(),
		Required:     1,
		SID:          "sid-abc",
		Resource:     "local",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kategori is not found")
	assert.Empty(t, res)
}

func TestCreateTemplatePertanyaanCommandHandler_BankSoalNotFound(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	repoKategori := &mockKategori.MockKategoriRepository{}
	repoBankSoal := &mockBankSoal.MockRepository{}

	uuidKategori := uuid.New()
	uuidBankSoal := uuid.New()

	repoKategori.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domainKategori.Kategori, error) {
		return &domainKategori.Kategori{ID: 10, UUID: uuidKategori}, nil
	}

	repoBankSoal.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domainBankSoal.BankSoal, error) {
		return nil, errors.New("not found")
	}

	handler := &application.CreateTemplatePertanyaanCommandHandler{
		Repo:         repo,
		RepoKategori: repoKategori,
		RepoBankSoal: repoBankSoal,
	}

	cmd := application.CreateTemplatePertanyaanCommand{
		UuidBankSoal: uuidBankSoal.String(),
		Pertanyaan:   "Apakah ini pertanyaan?",
		JenisPilihan: "pilihan_ganda",
		Bobot:        "5",
		UuidKategori: uuidKategori.String(),
		Required:     1,
		SID:          "sid-abc",
		Resource:     "local",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kategori is not found")
	assert.Empty(t, res)
}
