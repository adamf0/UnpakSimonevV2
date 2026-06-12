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

func TestGetBankSoalByUuidQueryHandler_Handle(t *testing.T) {
	validUUID := uuid.New()

	t.Run("invalid uuid", func(t *testing.T) {
		repo := &mock.MockRepository{}
		handler := &GetBankSoalByUuidQueryHandler{Repo: repo}
		cmd := GetBankSoalByUuidQuery{
			Uuid: "invalid-uuid",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("not found in db", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		handler := &GetBankSoalByUuidQueryHandler{Repo: repo}
		cmd := GetBankSoalByUuidQuery{
			Uuid: validUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("db general error", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return nil, errors.New("db error")
			},
		}
		handler := &GetBankSoalByUuidQueryHandler{Repo: repo}
		cmd := GetBankSoalByUuidQuery{
			Uuid: validUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success", func(t *testing.T) {
		expected := &domain.BankSoal{UUID: validUUID, Judul: "Success Soal"}
		repo := &mock.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				assert.Equal(t, validUUID, uid)
				return expected, nil
			},
		}
		handler := &GetBankSoalByUuidQueryHandler{Repo: repo}
		cmd := GetBankSoalByUuidQuery{
			Uuid: validUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}
