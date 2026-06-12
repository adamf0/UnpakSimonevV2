package application

import (
	"context"
	"errors"
	"testing"

	"UnpakSiamida/modules/banksoal/application/mock"

	"github.com/stretchr/testify/assert"
)

func TestSetupUuidBankSoalCommandHandler_Handle(t *testing.T) {
	t.Run("repo setup uuid error", func(t *testing.T) {
		repo := &mock.MockRepository{
			SetupUuidFunc: func(ctx context.Context) error {
				return errors.New("setup failed")
			},
		}
		handler := &SetupUuidBankSoalCommandHandler{Repo: repo}
		cmd := SetupUuidBankSoalCommand{}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("success", func(t *testing.T) {
		repo := &mock.MockRepository{
			SetupUuidFunc: func(ctx context.Context) error {
				return nil
			},
		}
		handler := &SetupUuidBankSoalCommandHandler{Repo: repo}
		cmd := SetupUuidBankSoalCommand{}
		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, "berhasil setup uuid pada data", res)
	})
}
