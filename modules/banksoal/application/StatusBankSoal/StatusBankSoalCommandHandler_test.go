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

func TestStatusBankSoalCommandHandler_Handle(t *testing.T) {
	validUUID := uuid.New()

	t.Run("invalid uuid", func(t *testing.T) {
		repo := &mock.MockRepository{}
		handler := &StatusBankSoalCommandHandler{Repo: repo}
		cmd := StatusBankSoalCommand{
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
		handler := &StatusBankSoalCommandHandler{Repo: repo}
		cmd := StatusBankSoalCommand{
			Uuid: validUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("get by uuid general error", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return nil, errors.New("db error")
			},
		}
		handler := &StatusBankSoalCommandHandler{Repo: repo}
		cmd := StatusBankSoalCommand{
			Uuid: validUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("invalid status", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{UUID: validUUID, Status: "draf"}, nil
			},
		}
		handler := &StatusBankSoalCommandHandler{Repo: repo}
		cmd := StatusBankSoalCommand{
			Uuid:   validUUID.String(),
			Status: "invalid-status",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("repo update error", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{UUID: validUUID, Status: "draf"}, nil
			},
			UpdateFunc: func(ctx context.Context, banksoal *domain.BankSoal) error {
				return errors.New("update error")
			},
		}
		handler := &StatusBankSoalCommandHandler{Repo: repo}
		cmd := StatusBankSoalCommand{
			Uuid:   validUUID.String(),
			Status: "active",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("success", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{UUID: validUUID, Status: "draf"}, nil
			},
			UpdateFunc: func(ctx context.Context, banksoal *domain.BankSoal) error {
				assert.Equal(t, "active", banksoal.Status)
				return nil
			},
		}
		handler := &StatusBankSoalCommandHandler{Repo: repo}
		cmd := StatusBankSoalCommand{
			Uuid:   validUUID.String(),
			Status: "active",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, validUUID.String(), res)
	})
}
