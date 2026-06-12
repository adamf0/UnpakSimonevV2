package application

import (
	"context"
	"errors"
	"testing"

	"UnpakSiamida/common/helper"
	mockrepo "UnpakSiamida/modules/account/application/mock"
	domainaccount "UnpakSiamida/modules/account/domain"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestLoginCommandHandler_Handle(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expectedUser := &domainaccount.AccountDefault{
			ID:       "user-123",
			Resource: helper.StrPtr("resource-val"),
			CodeCtx:  helper.StrPtr("code-ctx-val"),
		}

		repo := &mockrepo.MockAccountRepository{
			AuthFunc: func(ctx context.Context, username, password string) (*domainaccount.AccountDefault, error) {
				assert.Equal(t, "testuser", username)
				assert.Equal(t, "testpass", password)
				return expectedUser, nil
			},
		}

		handler := &LoginCommandHandler{
			Repo: repo,
		}

		cmd := LoginCommand{
			Username: "testuser",
			Password: "testpass",
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.AccessToken)
		assert.NotEmpty(t, res.RefreshToken)
		assert.Equal(t, "user-123", res.UserID)
		assert.Equal(t, "resource-val", res.Resource)
		assert.Equal(t, "code-ctx-val", res.CodeCtx)
	})

	t.Run("Failure_InvalidCredential_NotFound", func(t *testing.T) {
		repo := &mockrepo.MockAccountRepository{
			AuthFunc: func(ctx context.Context, username, password string) (*domainaccount.AccountDefault, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		handler := &LoginCommandHandler{
			Repo: repo,
		}

		cmd := LoginCommand{
			Username: "testuser",
			Password: "testpass",
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainaccount.InvalidCredential().Error(), err.Error())
		assert.Nil(t, res)
	})

	t.Run("Failure_RepoError", func(t *testing.T) {
		expectedErr := errors.New("auth DB error")
		repo := &mockrepo.MockAccountRepository{
			AuthFunc: func(ctx context.Context, username, password string) (*domainaccount.AccountDefault, error) {
				return nil, expectedErr
			},
		}

		handler := &LoginCommandHandler{
			Repo: repo,
		}

		cmd := LoginCommand{
			Username: "testuser",
			Password: "testpass",
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, res)
	})
}
