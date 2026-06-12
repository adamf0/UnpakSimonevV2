package application

import (
	"UnpakSiamida/modules/kategori/application/mock"
	domainkategori "UnpakSiamida/modules/kategori/domain"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteKategoriCommandHandler_Handle(t *testing.T) {
	targetUUID := uuid.New()
	targetUUIDStr := targetUUID.String()

	existingKategori := &domainkategori.Kategori{
		ID:           123,
		UUID:         targetUUID,
		NamaKategori: "Some Category",
	}

	t.Run("Success case - soft delete", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				assert.Equal(t, targetUUID, uid)
				return existingKategori, nil
			},
			UpdateFunc: func(ctx context.Context, kategori *domainkategori.Kategori) error {
				assert.NotNil(t, kategori.DeletedAt)
				return nil
			},
		}

		handler := &DeleteKategoriCommandHandler{Repo: repo}
		cmd := DeleteKategoriCommand{
			Uuid: targetUUIDStr,
			Mode: "soft_delete",
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, targetUUIDStr, res)
	})

	t.Run("Success case - hard delete", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				assert.Equal(t, targetUUID, uid)
				return existingKategori, nil
			},
			DeleteFunc: func(ctx context.Context, uid uuid.UUID) error {
				assert.Equal(t, targetUUID, uid)
				return nil
			},
		}

		handler := &DeleteKategoriCommandHandler{Repo: repo}
		cmd := DeleteKategoriCommand{
			Uuid: targetUUIDStr,
			Mode: "hard_delete",
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, targetUUIDStr, res)
	})

	t.Run("Failure case - invalid UUID format", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{}
		handler := &DeleteKategoriCommandHandler{Repo: repo}
		cmd := DeleteKategoriCommand{
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

		handler := &DeleteKategoriCommandHandler{Repo: repo}
		cmd := DeleteKategoriCommand{
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

		handler := &DeleteKategoriCommandHandler{Repo: repo}
		cmd := DeleteKategoriCommand{
			Uuid: targetUUIDStr,
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})

	t.Run("Failure case - db error on Delete (hard delete)", func(t *testing.T) {
		dbErr := errors.New("delete error")
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				return existingKategori, nil
			},
			DeleteFunc: func(ctx context.Context, uid uuid.UUID) error {
				return dbErr
			},
		}

		handler := &DeleteKategoriCommandHandler{Repo: repo}
		cmd := DeleteKategoriCommand{
			Uuid: targetUUIDStr,
			Mode: "hard_delete",
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})

	t.Run("Failure case - db error on Update (soft delete)", func(t *testing.T) {
		dbErr := errors.New("update error")
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				return existingKategori, nil
			},
			UpdateFunc: func(ctx context.Context, kategori *domainkategori.Kategori) error {
				return dbErr
			},
		}

		handler := &DeleteKategoriCommandHandler{Repo: repo}
		cmd := DeleteKategoriCommand{
			Uuid: targetUUIDStr,
			Mode: "soft_delete",
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})
}
