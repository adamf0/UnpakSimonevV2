package application

import (
	"context"
	"errors"
	"testing"

	mockkuesioner "UnpakSiamida/modules/kuesioner/application/mock"
	domainkuesioner "UnpakSiamida/modules/kuesioner/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetKuesionerByUuidQueryHandler(t *testing.T) {
	validUUID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		expectedKuesioner := &domainkuesioner.Kuesioner{
			UUID: validUUID,
		}
		repo := &mockkuesioner.MockKuesionerRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.Kuesioner, error) {
				assert.Equal(t, validUUID, uid)
				return expectedKuesioner, nil
			},
		}

		handler := &GetKuesionerByUuidQueryHandler{
			Repo: repo,
		}

		q := GetKuesionerByUuidQuery{
			Uuid: validUUID.String(),
		}

		res, err := handler.Handle(context.Background(), q)
		assert.NoError(t, err)
		assert.Equal(t, expectedKuesioner, res)
	})

	t.Run("Fail Parse UUID Error", func(t *testing.T) {
		handler := &GetKuesionerByUuidQueryHandler{}
		q := GetKuesionerByUuidQuery{
			Uuid: "invalid-uuid",
		}
		_, err := handler.Handle(context.Background(), q)
		assert.Error(t, err)
		assert.Equal(t, domainkuesioner.NotFound("invalid-uuid"), err)
	})

	t.Run("Fail Record Not Found", func(t *testing.T) {
		repo := &mockkuesioner.MockKuesionerRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.Kuesioner, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		handler := &GetKuesionerByUuidQueryHandler{
			Repo: repo,
		}

		q := GetKuesionerByUuidQuery{
			Uuid: validUUID.String(),
		}

		_, err := handler.Handle(context.Background(), q)
		assert.Error(t, err)
		assert.Equal(t, domainkuesioner.NotFound(validUUID.String()), err)
	})

	t.Run("Fail Other DB Error", func(t *testing.T) {
		dbErr := errors.New("db error")
		repo := &mockkuesioner.MockKuesionerRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.Kuesioner, error) {
				return nil, dbErr
			},
		}

		handler := &GetKuesionerByUuidQueryHandler{
			Repo: repo,
		}

		q := GetKuesionerByUuidQuery{
			Uuid: validUUID.String(),
		}

		_, err := handler.Handle(context.Background(), q)
		assert.ErrorIs(t, err, dbErr)
	})
}
