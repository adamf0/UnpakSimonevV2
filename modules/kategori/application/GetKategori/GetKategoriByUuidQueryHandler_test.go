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

func TestGetKategoriByUuidQueryHandler_Handle(t *testing.T) {
	targetUUID := uuid.New()
	targetUUIDStr := targetUUID.String()

	existingKategori := &domainKategori.Kategori{
		ID:           123,
		UUID:         targetUUID,
		NamaKategori: "Test Category",
	}

	t.Run("Success case", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainKategori.Kategori, error) {
				assert.Equal(t, targetUUID, uid)
				return existingKategori, nil
			},
		}

		handler := &GetKategoriByUuidQueryHandler{Repo: repo}
		q := GetKategoriByUuidQuery{
			Uuid: targetUUIDStr,
		}

		res, err := handler.Handle(context.Background(), q)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, targetUUID, res.UUID)
		assert.Equal(t, "Test Category", res.NamaKategori)
	})

	t.Run("Failure case - invalid UUID format", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{}
		handler := &GetKategoriByUuidQueryHandler{Repo: repo}
		q := GetKategoriByUuidQuery{
			Uuid: "invalid-uuid",
		}

		_, err := handler.Handle(context.Background(), q)
		assert.Error(t, err)
		assert.Equal(t, domainKategori.NotFound("invalid-uuid"), err)
	})

	t.Run("Failure case - category not found", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainKategori.Kategori, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		handler := &GetKategoriByUuidQueryHandler{Repo: repo}
		q := GetKategoriByUuidQuery{
			Uuid: targetUUIDStr,
		}

		_, err := handler.Handle(context.Background(), q)
		assert.Error(t, err)
		assert.Equal(t, domainKategori.NotFound(targetUUIDStr), err)
	})

	t.Run("Failure case - database error", func(t *testing.T) {
		dbErr := errors.New("db error")
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainKategori.Kategori, error) {
				return nil, dbErr
			},
		}

		handler := &GetKategoriByUuidQueryHandler{Repo: repo}
		q := GetKategoriByUuidQuery{
			Uuid: targetUUIDStr,
		}

		_, err := handler.Handle(context.Background(), q)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})
}
