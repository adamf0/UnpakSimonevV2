package application

import (
	"context"
	"errors"
	"testing"

	mockkuesioner "UnpakSiamida/modules/kuesioner/application/mock"
	domainkuesioner "UnpakSiamida/modules/kuesioner/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteKuesionerCommandHandler(t *testing.T) {
	validUUID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		repo := &mockkuesioner.MockKuesionerRepository{
			DeleteFunc: func(ctx context.Context, uid uuid.UUID) error {
				assert.Equal(t, validUUID, uid)
				return nil
			},
		}
		handler := &DeleteKuesionerCommandHandler{
			Repo: repo,
		}
		cmd := DeleteKuesionerCommand{
			Uuid: validUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, validUUID.String(), res)
	})

	t.Run("Fail Invalid UUID", func(t *testing.T) {
		handler := &DeleteKuesionerCommandHandler{}
		cmd := DeleteKuesionerCommand{
			Uuid: "invalid-uuid",
		}
		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainkuesioner.InvalidUuid(), err)
	})

	t.Run("Fail Repo Delete Error", func(t *testing.T) {
		deleteErr := errors.New("delete error")
		repo := &mockkuesioner.MockKuesionerRepository{
			DeleteFunc: func(ctx context.Context, uid uuid.UUID) error {
				return deleteErr
			},
		}
		handler := &DeleteKuesionerCommandHandler{
			Repo: repo,
		}
		cmd := DeleteKuesionerCommand{
			Uuid: validUUID.String(),
		}
		_, err := handler.Handle(context.Background(), cmd)
		assert.ErrorIs(t, err, deleteErr)
	})
}
