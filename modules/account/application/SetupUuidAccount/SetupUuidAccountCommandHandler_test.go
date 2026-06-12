package application

import (
	"context"
	"errors"
	"testing"

	mockrepo "UnpakSiamida/modules/account/application/mock"

	"github.com/stretchr/testify/assert"
)

func TestSetupUuidAccountCommandHandler_Handle(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := &mockrepo.MockAccountRepository{
			SetupUuidFunc: func(ctx context.Context) error {
				return nil
			},
		}

		handler := &SetupUuidAccountCommandHandler{
			Repo: repo,
		}

		res, err := handler.Handle(context.Background(), SetupUuidAccountCommand{})
		assert.NoError(t, err)
		assert.Equal(t, "berhasil setup uuid pada data", res)
	})

	t.Run("Failure", func(t *testing.T) {
		expectedErr := errors.New("setup uuid failure")
		repo := &mockrepo.MockAccountRepository{
			SetupUuidFunc: func(ctx context.Context) error {
				return expectedErr
			},
		}

		handler := &SetupUuidAccountCommandHandler{
			Repo: repo,
		}

		res, err := handler.Handle(context.Background(), SetupUuidAccountCommand{})
		assert.ErrorIs(t, err, expectedErr)
		assert.Empty(t, res)
	})
}
