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

func TestGetKuesionerDefaultByUuidQueryHandler(t *testing.T) {
	validUUID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		expectedKuesioner := &domainkuesioner.KuesionerDefault{
			UUID:  validUUID,
			Judul: "Kuesioner Default A",
		}
		repo := &mockkuesioner.MockKuesionerRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.KuesionerDefault, error) {
				assert.Equal(t, validUUID, uid)
				return expectedKuesioner, nil
			},
		}

		handler := &GetKuesionerDefaultByUuidQueryHandler{
			Repo: repo,
		}

		q := GetKuesionerDefaultByUuidQuery{
			Uuid: validUUID.String(),
		}

		res, err := handler.Handle(context.Background(), q)
		assert.NoError(t, err)
		assert.Equal(t, expectedKuesioner, res)
	})

	t.Run("Fail Parse UUID Error", func(t *testing.T) {
		handler := &GetKuesionerDefaultByUuidQueryHandler{}
		q := GetKuesionerDefaultByUuidQuery{
			Uuid: "invalid-uuid",
		}
		_, err := handler.Handle(context.Background(), q)
		assert.Error(t, err)
		assert.Equal(t, domainkuesioner.NotFound("invalid-uuid"), err)
	})

	t.Run("Fail Record Not Found", func(t *testing.T) {
		repo := &mockkuesioner.MockKuesionerRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.KuesionerDefault, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		handler := &GetKuesionerDefaultByUuidQueryHandler{
			Repo: repo,
		}

		q := GetKuesionerDefaultByUuidQuery{
			Uuid: validUUID.String(),
		}

		_, err := handler.Handle(context.Background(), q)
		assert.Error(t, err)
		assert.Equal(t, domainkuesioner.NotFound(validUUID.String()), err)
	})

	t.Run("Fail Other DB Error", func(t *testing.T) {
		dbErr := errors.New("db error")
		repo := &mockkuesioner.MockKuesionerRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.KuesionerDefault, error) {
				return nil, dbErr
			},
		}

		handler := &GetKuesionerDefaultByUuidQueryHandler{
			Repo: repo,
		}

		q := GetKuesionerDefaultByUuidQuery{
			Uuid: validUUID.String(),
		}

		_, err := handler.Handle(context.Background(), q)
		assert.ErrorIs(t, err, dbErr)
	})
}
