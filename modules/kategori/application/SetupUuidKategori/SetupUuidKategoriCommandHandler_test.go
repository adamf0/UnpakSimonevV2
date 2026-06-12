package application

import (
	"UnpakSiamida/modules/kategori/application/mock"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupUuidKategoriCommandHandler_Handle(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{
			SetupUuidFunc: func(ctx context.Context) error {
				return nil
			},
		}

		handler := &SetupUuidKategoriCommandHandler{Repo: repo}
		cmd := SetupUuidKategoriCommand{}

		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, "berhasil setup uuid pada data", res)
	})

	t.Run("Failure case - SetupUuid error", func(t *testing.T) {
		dbErr := errors.New("db error")
		repo := &mock.MockKategoriRepository{
			SetupUuidFunc: func(ctx context.Context) error {
				return dbErr
			},
		}

		handler := &SetupUuidKategoriCommandHandler{Repo: repo}
		cmd := SetupUuidKategoriCommand{}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})
}
