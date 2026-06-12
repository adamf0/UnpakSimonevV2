package application_test

import (
	"context"
	"sync"
	"testing"

	copyjawaban "UnpakSiamida/modules/templatejawaban/application/CopyTemplateJawaban"
	"UnpakSiamida/modules/templatepertanyaan/application/CopyTemplatePertanyaan"
	"UnpakSiamida/modules/templatepertanyaan/application/mock"
	"UnpakSiamida/modules/templatepertanyaan/domain"

	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type mockCopyTemplateJawabanHandler struct {
	HandleFunc func(ctx context.Context, cmd copyjawaban.CopyTemplateJawabanCommand) (string, error)
}

func (h *mockCopyTemplateJawabanHandler) Handle(ctx context.Context, cmd copyjawaban.CopyTemplateJawabanCommand) (string, error) {
	if h.HandleFunc != nil {
		return h.HandleFunc(ctx, cmd)
	}
	return "success", nil
}

var registerOnce sync.Once
var mockCopyJawaban mockCopyTemplateJawabanHandler

func initMediator() {
	registerOnce.Do(func() {
		_ = mediatr.RegisterRequestHandler[copyjawaban.CopyTemplateJawabanCommand, string](&mockCopyJawaban)
	})
}

func TestCopyTemplatePertanyaanCommandHandler_Success(t *testing.T) {
	initMediator()
	mockCopyJawaban.HandleFunc = func(ctx context.Context, cmd copyjawaban.CopyTemplateJawabanCommand) (string, error) {
		return "success-id", nil
	}

	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	existingTp := &domain.TemplatePertanyaan{
		ID:           1,
		UUID:         uid,
		IdBankSoal:   2,
		Pertanyaan:   "Pertanyaan lama",
		JenisPilihan: "pilihan_ganda",
		Bobot:        1,
		Required:     1,
		Status:       "draf",
	}

	repo.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaan, error) {
		assert.Equal(t, uid, id)
		return existingTp, nil
	}

	repo.CountCopyFunc = func(ctx context.Context, judul string) (int, error) {
		assert.Equal(t, "Pertanyaan lama", judul)
		return 0, nil
	}

	repo.CreateFunc = func(ctx context.Context, tp *domain.TemplatePertanyaan) error {
		assert.Contains(t, tp.Pertanyaan, "salin - Pertanyaan lama")
		assert.Equal(t, uint(2), tp.IdBankSoal)
		return nil
	}

	handler := &application.CopyTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.CopyTemplatePertanyaanCommand{
		Uuid:     uid.String(),
		SID:      "sid-123",
		Resource: "local",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestCopyTemplatePertanyaanCommandHandler_FailGetByUuid(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	repo.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaan, error) {
		return nil, gorm.ErrRecordNotFound
	}

	handler := &application.CopyTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.CopyTemplatePertanyaanCommand{
		Uuid:     uid.String(),
		SID:      "sid-123",
		Resource: "local",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Empty(t, res)
}

func TestCopyTemplatePertanyaanCommandHandler_InvalidUuid(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}

	handler := &application.CopyTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.CopyTemplatePertanyaanCommand{
		Uuid:     "invalid-uuid",
		SID:      "sid-123",
		Resource: "local",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "uuid is invalid")
	assert.Empty(t, res)
}

func TestCopyTemplatePertanyaanResultCommandHandler(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	expectedMap := map[uint]uint{1: 2}

	repo.CopyByBankSoalFunc = func(ctx context.Context, tx *gorm.DB, sourceBankSoalID uint, targetBankSoalID uint, resource string, sid string) (map[uint]uint, error) {
		assert.Equal(t, uint(1), sourceBankSoalID)
		assert.Equal(t, uint(2), targetBankSoalID)
		assert.Equal(t, "local", resource)
		assert.Equal(t, "sid-123", sid)
		return expectedMap, nil
	}

	handler := &application.CopyTemplatePertanyaanResultCommandHandler{
		Repo: repo,
	}

	cmd := application.CopyTemplatePertanyaanResultCommand{
		Tx:               nil,
		SourceBankSoalID: 1,
		TargetBankSoalID: 2,
		Resource:         "local",
		Sid:              "sid-123",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, expectedMap, res)
}
