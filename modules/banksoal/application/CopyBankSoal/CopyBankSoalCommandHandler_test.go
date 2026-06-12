package application

import (
	"context"
	"errors"
	"sync"
	"testing"

	"UnpakSiamida/modules/banksoal/application/mock"
	"UnpakSiamida/modules/banksoal/domain"
	copyJawaban "UnpakSiamida/modules/templatejawaban/application/CopyTemplateJawaban"
	copy "UnpakSiamida/modules/templatepertanyaan/application/CopyTemplatePertanyaan"

	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	mockCopyPertanyaanFunc func(ctx context.Context, cmd copy.CopyTemplatePertanyaanResultCommand) (map[uint]uint, error)
	mockCopyJawabanFunc    func(ctx context.Context, cmd copyJawaban.CopyTemplateJawabanCommand) (string, error)
	registerOnce           sync.Once
)

type mockCopyPertanyaanHandler struct{}

func (h *mockCopyPertanyaanHandler) Handle(ctx context.Context, cmd copy.CopyTemplatePertanyaanResultCommand) (map[uint]uint, error) {
	if mockCopyPertanyaanFunc != nil {
		return mockCopyPertanyaanFunc(ctx, cmd)
	}
	return map[uint]uint{1: 2}, nil
}

type mockCopyJawabanHandler struct{}

func (h *mockCopyJawabanHandler) Handle(ctx context.Context, cmd copyJawaban.CopyTemplateJawabanCommand) (string, error) {
	if mockCopyJawabanFunc != nil {
		return mockCopyJawabanFunc(ctx, cmd)
	}
	return "success", nil
}

func setupMediatrMocks() {
	registerOnce.Do(func() {
		_ = mediatr.RegisterRequestHandler[copy.CopyTemplatePertanyaanResultCommand, map[uint]uint](&mockCopyPertanyaanHandler{})
		_ = mediatr.RegisterRequestHandler[copyJawaban.CopyTemplateJawabanCommand, string](&mockCopyJawabanHandler{})
	})
}

func TestCopyBankSoalCommandHandler_Handle(t *testing.T) {
	setupMediatrMocks()

	validUUID := uuid.New()
	targetUUID := uuid.New()

	t.Run("invalid uuid format", func(t *testing.T) {
		repo := &mock.MockRepository{}
		handler := &CopyBankSoalCommandHandler{Repo: repo}
		cmd := CopyBankSoalCommand{
			Uuid:     "invalid-uuid",
			Resource: "local",
			SID:      "sid-123",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("bank soal not found", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		handler := &CopyBankSoalCommandHandler{Repo: repo}
		cmd := CopyBankSoalCommand{
			Uuid:     validUUID.String(),
			Resource: "local",
			SID:      "sid-123",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("count copy repo error", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{
					ID:   1,
					UUID: validUUID,
				}, nil
			},
			CountCopyFunc: func(ctx context.Context, judul string) (int, error) {
				return 0, errors.New("database error")
			},
		}
		handler := &CopyBankSoalCommandHandler{Repo: repo}
		cmd := CopyBankSoalCommand{
			Uuid:     validUUID.String(),
			Resource: "local",
			SID:      "sid-123",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("success", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{
					ID:    1,
					UUID:  validUUID,
					Judul: "Existing Soal",
				}, nil
			},
			CountCopyFunc: func(ctx context.Context, judul string) (int, error) {
				return 1, nil
			},
			CreateFunc: func(ctx context.Context, banksoal *domain.BankSoal) error {
				banksoal.ID = 2
				banksoal.UUID = targetUUID
				return nil
			},
		}
		mockCopyPertanyaanFunc = func(ctx context.Context, cmd copy.CopyTemplatePertanyaanResultCommand) (map[uint]uint, error) {
			return map[uint]uint{10: 20}, nil
		}
		mockCopyJawabanFunc = func(ctx context.Context, cmd copyJawaban.CopyTemplateJawabanCommand) (string, error) {
			return "success", nil
		}

		handler := &CopyBankSoalCommandHandler{Repo: repo}
		cmd := CopyBankSoalCommand{
			Uuid:     validUUID.String(),
			Resource: "local",
			SID:      "sid-123",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
