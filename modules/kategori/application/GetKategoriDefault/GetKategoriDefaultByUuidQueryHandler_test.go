package application

import (
	"UnpakSiamida/modules/kategori/application/mock"
	domainKategori "UnpakSiamida/modules/kategori/domain"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetKategoriDefaultByUuidQueryHandler_Handle(t *testing.T) {
	targetUUID := uuid.New()
	targetUUIDStr := targetUUID.String()

	existingKategoriDefault := &domainKategori.KategoriDefault{
		ID:           123,
		UUID:         targetUUID,
		NamaKategori: "Test Category Default",
	}

	t.Run("Success case", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainKategori.KategoriDefault, error) {
				assert.Equal(t, targetUUID, uid)
				return existingKategoriDefault, nil
			},
		}

		handler := &GetKategoriDefaultByUuidQueryHandler{Repo: repo}
		q := GetKategoriDefaultByUuidQuery{
			Uuid: targetUUIDStr,
		}

		res, err := handler.Handle(context.Background(), q)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, targetUUID, res.UUID)
		assert.Equal(t, "Test Category Default", res.NamaKategori)
	})

	t.Run("Failure case - invalid UUID format", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{}
		handler := &GetKategoriDefaultByUuidQueryHandler{Repo: repo}
		q := GetKategoriDefaultByUuidQuery{
			Uuid: "invalid-uuid",
		}

		_, err := handler.Handle(context.Background(), q)
		assert.Error(t, err)
		assert.Equal(t, domainKategori.NotFound("invalid-uuid"), err)
	})

	t.Run("Failure case - category default not found", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainKategori.KategoriDefault, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		handler := &GetKategoriDefaultByUuidQueryHandler{Repo: repo}
		q := GetKategoriDefaultByUuidQuery{
			Uuid: targetUUIDStr,
		}

		_, err := handler.Handle(context.Background(), q)
		assert.Error(t, err)
		assert.Equal(t, domainKategori.NotFound(targetUUIDStr), err)
	})

	t.Run("Failure case - database error", func(t *testing.T) {
		dbErr := errors.New("db error")
		repo := &mock.MockKategoriRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainKategori.KategoriDefault, error) {
				return nil, dbErr
			},
		}

		handler := &GetKategoriDefaultByUuidQueryHandler{Repo: repo}
		q := GetKategoriDefaultByUuidQuery{
			Uuid: targetUUIDStr,
		}

		_, err := handler.Handle(context.Background(), q)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})
}
