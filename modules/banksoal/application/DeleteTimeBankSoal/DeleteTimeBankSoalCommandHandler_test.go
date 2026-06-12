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

func TestDeleteTimeBankSoalCommandHandler_Handle(t *testing.T) {
	validUUID := uuid.New()

	t.Run("invalid uuid", func(t *testing.T) {
		repo := &mock.MockRepository{}
		handler := &DeleteTimeBankSoalCommandHandler{Repo: repo}
		cmd := DeleteTimeBankSoalCommand{
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
		handler := &DeleteTimeBankSoalCommandHandler{Repo: repo}
		cmd := DeleteTimeBankSoalCommand{
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
		handler := &DeleteTimeBankSoalCommandHandler{Repo: repo}
		cmd := DeleteTimeBankSoalCommand{
			Uuid: validUUID.String(),
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
				assert.Nil(t, banksoal.TanggalMulai)
				assert.Nil(t, banksoal.TanggalAkhir)
				return nil
			},
		}
		handler := &DeleteTimeBankSoalCommandHandler{Repo: repo}
		cmd := DeleteTimeBankSoalCommand{
			Uuid: validUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, validUUID.String(), res)
	})

	t.Run("update repo error", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{UUID: validUUID}, nil
			},
			UpdateFunc: func(ctx context.Context, banksoal *domain.BankSoal) error {
				return errors.New("update failed")
			},
		}
		handler := &DeleteTimeBankSoalCommandHandler{Repo: repo}
		cmd := DeleteTimeBankSoalCommand{
			Uuid: validUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestDeleteTimeExtBankSoalCommandHandler_Handle(t *testing.T) {
	validExtUUID := uuid.New()
	validBankSoalUUID := uuid.New()

	t.Run("invalid ext uuid", func(t *testing.T) {
		repo := &mock.MockRepository{}
		handler := &DeleteTimeExtBankSoalCommandHandler{Repo: repo}
		cmd := DeleteTimeExtBankSoalCommand{
			Uuid:         "invalid-uuid",
			UuidBankSoal: validBankSoalUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("invalid bank soal uuid", func(t *testing.T) {
		repo := &mock.MockRepository{}
		handler := &DeleteTimeExtBankSoalCommandHandler{Repo: repo}
		cmd := DeleteTimeExtBankSoalCommand{
			Uuid:         validExtUUID.String(),
			UuidBankSoal: "invalid-uuid",
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
		handler := &DeleteTimeExtBankSoalCommandHandler{Repo: repo}
		cmd := DeleteTimeExtBankSoalCommand{
			Uuid:         validExtUUID.String(),
			UuidBankSoal: validBankSoalUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("get bank soal general error", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return nil, errors.New("db error")
			},
		}
		handler := &DeleteTimeExtBankSoalCommandHandler{Repo: repo}
		cmd := DeleteTimeExtBankSoalCommand{
			Uuid:         validExtUUID.String(),
			UuidBankSoal: validBankSoalUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("delete ext repo error", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{ID: 10, UUID: validBankSoalUUID}, nil
			},
			DeleteExtFunc: func(ctx context.Context, uid uuid.UUID, idbanksoal uint) error {
				assert.Equal(t, validExtUUID, uid)
				assert.Equal(t, uint(10), idbanksoal)
				return errors.New("delete ext failed")
			},
		}
		handler := &DeleteTimeExtBankSoalCommandHandler{Repo: repo}
		cmd := DeleteTimeExtBankSoalCommand{
			Uuid:         validExtUUID.String(),
			UuidBankSoal: validBankSoalUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("success", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{ID: 10, UUID: validBankSoalUUID}, nil
			},
			DeleteExtFunc: func(ctx context.Context, uid uuid.UUID, idbanksoal uint) error {
				assert.Equal(t, validExtUUID, uid)
				assert.Equal(t, uint(10), idbanksoal)
				return nil
			},
		}
		handler := &DeleteTimeExtBankSoalCommandHandler{Repo: repo}
		cmd := DeleteTimeExtBankSoalCommand{
			Uuid:         validExtUUID.String(),
			UuidBankSoal: validBankSoalUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, validExtUUID.String(), res)
	})
}
