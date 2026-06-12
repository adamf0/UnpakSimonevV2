package application

import (
	"context"
	"errors"
	"testing"

	"UnpakSiamida/modules/banksoal/application/mock"
	"UnpakSiamida/modules/banksoal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUpdateBankSoalCommandHandler_Handle(t *testing.T) {
	validUUID := uuid.New()

	t.Run("invalid uuid", func(t *testing.T) {
		repo := &mock.MockRepository{}
		handler := &UpdateBankSoalCommandHandler{Repo: repo}
		cmd := UpdateBankSoalCommand{
			Uuid: "invalid-uuid",
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
		handler := &UpdateBankSoalCommandHandler{Repo: repo}
		cmd := UpdateBankSoalCommand{
			Uuid: validUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("invalid owner / resource is not local", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{UUID: validUUID}, nil
			},
		}
		handler := &UpdateBankSoalCommandHandler{Repo: repo}
		cmd := UpdateBankSoalCommand{
			Uuid:     validUUID.String(),
			Resource: "not-local",
			SID:      "sid-123",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("repo update error", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{UUID: validUUID}, nil
			},
			UpdateFunc: func(ctx context.Context, banksoal *domain.BankSoal) error {
				return errors.New("update error")
			},
		}
		handler := &UpdateBankSoalCommandHandler{Repo: repo}
		cmd := UpdateBankSoalCommand{
			Uuid:     validUUID.String(),
			Judul:    "Updated Judul",
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
				return &domain.BankSoal{UUID: validUUID}, nil
			},
			UpdateFunc: func(ctx context.Context, banksoal *domain.BankSoal) error {
				assert.Equal(t, "Updated Judul", banksoal.Judul)
				return nil
			},
		}
		handler := &UpdateBankSoalCommandHandler{Repo: repo}
		cmd := UpdateBankSoalCommand{
			Uuid:     validUUID.String(),
			Judul:    "Updated Judul",
			Resource: "local",
			SID:      "sid-123",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, validUUID.String(), res)
	})
}
