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

func TestCopyKategoriCommandHandler_Handle(t *testing.T) {
	targetUUID := uuid.New()
	targetUUIDStr := targetUUID.String()

	existingKategori := &domainkategori.Kategori{
		ID:           123,
		UUID:         targetUUID,
		NamaKategori: "Some Category",
	}

	t.Run("Success case", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				assert.Equal(t, targetUUID, uid)
				return existingKategori, nil
			},
			CountCopyFunc: func(ctx context.Context, judul string) (int, error) {
				assert.Equal(t, "Some Category", judul)
				return 2, nil
			},
			CreateFunc: func(ctx context.Context, kategori *domainkategori.Kategori) error {
				assert.Equal(t, "salin (3) - Some Category", kategori.NamaKategori)
				return nil
			},
			RebuildFullTextFunc: func(ctx context.Context) error {
				return nil
			},
		}

		handler := &CopyKategoriCommandHandler{Repo: repo}
		cmd := CopyKategoriCommand{
			Uuid:     targetUUIDStr,
			Resource: "local",
			SID:      "sid",
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("Failure case - invalid UUID format", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{}
		handler := &CopyKategoriCommandHandler{Repo: repo}
		cmd := CopyKategoriCommand{
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

		handler := &CopyKategoriCommandHandler{Repo: repo}
		cmd := CopyKategoriCommand{
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

		handler := &CopyKategoriCommandHandler{Repo: repo}
		cmd := CopyKategoriCommand{
			Uuid: targetUUIDStr,
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})

	t.Run("Failure case - db error on CountCopy", func(t *testing.T) {
		dbErr := errors.New("count error")
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				return existingKategori, nil
			},
			CountCopyFunc: func(ctx context.Context, judul string) (int, error) {
				return 0, dbErr
			},
		}

		handler := &CopyKategoriCommandHandler{Repo: repo}
		cmd := CopyKategoriCommand{
			Uuid: targetUUIDStr,
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})

	t.Run("Failure case - BeginTx fails", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				return existingKategori, nil
			},
			CountCopyFunc: func(ctx context.Context, judul string) (int, error) {
				return 0, nil
			},
			BeginTxFunc: func(ctx context.Context) (*gorm.DB, error) {
				return nil, errors.New("begin tx error")
			},
		}

		handler := &CopyKategoriCommandHandler{Repo: repo}
		cmd := CopyKategoriCommand{
			Uuid:     targetUUIDStr,
			Resource: "local",
			SID:      "sid",
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
	})
}
