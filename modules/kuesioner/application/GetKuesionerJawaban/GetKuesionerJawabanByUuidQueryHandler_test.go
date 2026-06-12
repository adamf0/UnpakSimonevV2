package application

import (
	"context"
	"errors"
	"testing"

	mockkuesioner "UnpakSiamida/modules/kuesioner/application/mock"
	domainkuesioner "UnpakSiamida/modules/kuesioner/domain"

	"github.com/stretchr/testify/assert"
)

func TestGetKuesionerJawabanByUuidQueryQueryHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expectedJawaban := []domainkuesioner.KuesionerJawabanDefault{
			{
				ID:       1,
				FreeText: nil,
			},
		}
		repo := &mockkuesioner.MockKuesionerJawabanRepository{
			GetAllByKuesionerFunc: func(ctx context.Context, uuidkuesioner string) ([]domainkuesioner.KuesionerJawabanDefault, error) {
				assert.Equal(t, "kuesioner-uuid-123", uuidkuesioner)
				return expectedJawaban, nil
			},
		}

		handler := &GetKuesionerJawabanByUuidQueryQueryHandler{
			Repo: repo,
		}

		q := GetKuesionerJawabanByUuidQuery{
			Uuid: "kuesioner-uuid-123",
		}

		res, err := handler.Handle(context.Background(), q)
		assert.NoError(t, err)
		assert.Equal(t, expectedJawaban, res)
	})

	t.Run("Fail Repo Error", func(t *testing.T) {
		repoErr := errors.New("repo error")
		repo := &mockkuesioner.MockKuesionerJawabanRepository{
			GetAllByKuesionerFunc: func(ctx context.Context, uuidkuesioner string) ([]domainkuesioner.KuesionerJawabanDefault, error) {
				return nil, repoErr
			},
		}

		handler := &GetKuesionerJawabanByUuidQueryQueryHandler{
			Repo: repo,
		}

		q := GetKuesionerJawabanByUuidQuery{
			Uuid: "kuesioner-uuid-123",
		}

		_, err := handler.Handle(context.Background(), q)
		assert.ErrorIs(t, err, repoErr)
	})
}
