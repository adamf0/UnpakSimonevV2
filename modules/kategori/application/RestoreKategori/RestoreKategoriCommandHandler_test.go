package application

import (
	"UnpakSiamida/modules/kategori/application/mock"
	domainkategori "UnpakSiamida/modules/kategori/domain"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRestoreKategoriCommandHandler_Handle(t *testing.T) {
	targetUUID := uuid.New()
	targetUUIDStr := targetUUID.String()
	now := time.Now()

	existingKategori := &domainkategori.Kategori{
		ID:           123,
		UUID:         targetUUID,
		DeletedAt:    &now,
		NamaKategori: "Some Category",
	}

	t.Run("Success case", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				assert.Equal(t, targetUUID, uid)
				return existingKategori, nil
			},
			UpdateFunc: func(ctx context.Context, kategori *domainkategori.Kategori) error {
				assert.Nil(t, kategori.DeletedAt)
				return nil
			},
		}

		handler := &RestoreKategoriCommandHandler{Repo: repo}
		cmd := RestoreKategoriCommand{
			Uuid: targetUUIDStr,
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, targetUUIDStr, res)
	})

	t.Run("Failure case - invalid UUID format", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{}
		handler := &RestoreKategoriCommandHandler{Repo: repo}
		cmd := RestoreKategoriCommand{
			Uuid: "invalid-uuid",
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainkategori.InvalidUuid(), err)
	})

	t.Run("Failure case - category not found", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		handler := &RestoreKategoriCommandHandler{Repo: repo}
		cmd := RestoreKategoriCommand{
			Uuid: targetUUIDStr,
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainkategori.NotFound(targetUUIDStr), err)
	})

	t.Run("Failure case - db error on GetByUuid", func(t *testing.T) {
		dbErr := errors.New("db error")
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				return nil, dbErr
			},
		}

		handler := &RestoreKategoriCommandHandler{Repo: repo}
		cmd := RestoreKategoriCommand{
			Uuid: targetUUIDStr,
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})

	t.Run("Failure case - db error on Update", func(t *testing.T) {
		dbErr := errors.New("update error")
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				return existingKategori, nil
			},
			UpdateFunc: func(ctx context.Context, kategori *domainkategori.Kategori) error {
				return dbErr
			},
		}

		handler := &RestoreKategoriCommandHandler{Repo: repo}
		cmd := RestoreKategoriCommand{
			Uuid: targetUUIDStr,
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})
}
