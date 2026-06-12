package application

import (
	"context"
	"errors"
	"testing"

	mockrepo "UnpakSiamida/modules/account/application/mock"
	domainaccount "UnpakSiamida/modules/account/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetAccountQueryHandler_Handle(t *testing.T) {
	validUUID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		expectedAccount := &domainaccount.Account{
			UUID: validUUID,
		}

		repo := &mockrepo.MockAccountRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainaccount.Account, error) {
				assert.Equal(t, validUUID, uid)
				return expectedAccount, nil
			},
		}

		handler := &GetAccountQueryHandler{
			Repo: repo,
		}

		q := GetAccountQuery{
			Uuid: validUUID.String(),
		}

		res, err := handler.Handle(context.Background(), q)
		assert.NoError(t, err)
		assert.Equal(t, expectedAccount, res)
	})

	t.Run("Failure_InvalidUuid", func(t *testing.T) {
		handler := &GetAccountQueryHandler{
			Repo: &mockrepo.MockAccountRepository{},
		}

		q := GetAccountQuery{
			Uuid: "invalid-uuid",
		}

		res, err := handler.Handle(context.Background(), q)
		assert.Error(t, err)
		assert.Equal(t, domainaccount.NotFound("invalid-uuid").Error(), err.Error())
		assert.Nil(t, res)
	})

	t.Run("Failure_NotFound", func(t *testing.T) {
		repo := &mockrepo.MockAccountRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainaccount.Account, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		handler := &GetAccountQueryHandler{
			Repo: repo,
		}

		q := GetAccountQuery{
			Uuid: validUUID.String(),
		}

		res, err := handler.Handle(context.Background(), q)
		assert.Error(t, err)
		assert.Equal(t, domainaccount.NotFound(validUUID.String()).Error(), err.Error())
		assert.Nil(t, res)
	})

	t.Run("Failure_RepoError", func(t *testing.T) {
		expectedErr := errors.New("db error")
		repo := &mockrepo.MockAccountRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainaccount.Account, error) {
				return nil, expectedErr
			},
		}

		handler := &GetAccountQueryHandler{
			Repo: repo,
		}

		q := GetAccountQuery{
			Uuid: validUUID.String(),
		}

		res, err := handler.Handle(context.Background(), q)
		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, res)
	})
}
