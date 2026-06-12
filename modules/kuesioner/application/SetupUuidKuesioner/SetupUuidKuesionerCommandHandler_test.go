package application

import (
	"context"
	"errors"
	"testing"

	mockkuesioner "UnpakSiamida/modules/kuesioner/application/mock"

	"github.com/stretchr/testify/assert"
)

func TestSetupUuidKuesionerCommandHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := &mockkuesioner.MockKuesionerRepository{
			SetupUuidFunc: func(ctx context.Context) error {
				return nil
			},
		}

		handler := &SetupUuidKuesionerCommandHandler{
			Repo: repo,
		}

		res, err := handler.Handle(context.Background(), SetupUuidKuesionerCommand{})
		assert.NoError(t, err)
		assert.Equal(t, "berhasil setup uuid pada data", res)
	})

	t.Run("Fail Repo SetupUuid Error", func(t *testing.T) {
		setupErr := errors.New("setup uuid error")
		repo := &mockkuesioner.MockKuesionerRepository{
			SetupUuidFunc: func(ctx context.Context) error {
				return setupErr
			},
		}

		handler := &SetupUuidKuesionerCommandHandler{
			Repo: repo,
		}

		_, err := handler.Handle(context.Background(), SetupUuidKuesionerCommand{})
		assert.ErrorIs(t, err, setupErr)
	})
}
