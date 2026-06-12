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

func TestDeleteBankSoalCommandHandler_Handle(t *testing.T) {
	validUUID := uuid.New()

	t.Run("invalid uuid", func(t *testing.T) {
		repo := &mock.MockRepository{}
		handler := &DeleteBankSoalCommandHandler{Repo: repo}
		cmd := DeleteBankSoalCommand{
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
		handler := &DeleteBankSoalCommandHandler{Repo: repo}
		cmd := DeleteBankSoalCommand{
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
		handler := &DeleteBankSoalCommandHandler{Repo: repo}
		cmd := DeleteBankSoalCommand{
			Uuid: validUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("hard delete success", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{UUID: validUUID}, nil
			},
			DeleteFunc: func(ctx context.Context, uid uuid.UUID) error {
				assert.Equal(t, validUUID, uid)
				return nil
			},
		}
		handler := &DeleteBankSoalCommandHandler{Repo: repo}
		cmd := DeleteBankSoalCommand{
			Uuid: validUUID.String(),
			Mode: "hard_delete",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, validUUID.String(), res)
	})

	t.Run("hard delete repo error", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{UUID: validUUID}, nil
			},
			DeleteFunc: func(ctx context.Context, uid uuid.UUID) error {
				return errors.New("delete failed")
			},
		}
		handler := &DeleteBankSoalCommandHandler{Repo: repo}
		cmd := DeleteBankSoalCommand{
			Uuid: validUUID.String(),
			Mode: "hard_delete",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("soft delete success", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{UUID: validUUID}, nil
			},
			UpdateFunc: func(ctx context.Context, banksoal *domain.BankSoal) error {
				assert.NotNil(t, banksoal.DeletedAt)
				return nil
			},
		}
		handler := &DeleteBankSoalCommandHandler{Repo: repo}
		cmd := DeleteBankSoalCommand{
			Uuid: validUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, validUUID.String(), res)
	})

	t.Run("soft delete update error", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{UUID: validUUID}, nil
			},
			UpdateFunc: func(ctx context.Context, banksoal *domain.BankSoal) error {
				return errors.New("update failed")
			},
		}
		handler := &DeleteBankSoalCommandHandler{Repo: repo}
		cmd := DeleteBankSoalCommand{
			Uuid: validUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}
