package application

import (
	"UnpakSiamida/modules/kategori/application/mock"
	domainkategori "UnpakSiamida/modules/kategori/domain"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateKategoriCommandHandler_Handle(t *testing.T) {
	t.Run("Success case - without parent", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{
			CreateFunc: func(ctx context.Context, kategori *domainkategori.Kategori) error {
				return nil
			},
			RebuildFullTextFunc: func(ctx context.Context) error {
				return nil
			},
		}

		handler := &CreateKategoriCommandHandler{Repo: repo}
		cmd := CreateKategoriCommand{
			NamaKategori: "Test Category",
			Resource:     "local",
			SID:          "sid",
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		_, parseErr := uuid.Parse(res)
		assert.NoError(t, parseErr)
	})

	t.Run("Success case - with valid parent", func(t *testing.T) {
		parentUUID := uuid.New()
		parentKategori := &domainkategori.Kategori{
			ID:           123,
			UUID:         parentUUID,
			NamaKategori: "Parent Category",
		}

		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				assert.Equal(t, parentUUID, uid)
				return parentKategori, nil
			},
			CreateFunc: func(ctx context.Context, kategori *domainkategori.Kategori) error {
				assert.Equal(t, uint(123), *kategori.SubKategori)
				return nil
			},
			RebuildFullTextFunc: func(ctx context.Context) error {
				return nil
			},
		}

		handler := &CreateKategoriCommandHandler{Repo: repo}
		parentUUIDStr := parentUUID.String()
		cmd := CreateKategoriCommand{
			NamaKategori: "Child Category",
			SubKategori:  &parentUUIDStr,
			Resource:     "local",
			SID:          "sid",
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("Failure case - invalid parent UUID format", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{}
		handler := &CreateKategoriCommandHandler{Repo: repo}
		badUUID := "not-a-uuid"
		cmd := CreateKategoriCommand{
			NamaKategori: "Child Category",
			SubKategori:  &badUUID,
			Resource:     "local",
			SID:          "sid",
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainkategori.InvalidUuid(), err)
	})

	t.Run("Failure case - parent not found", func(t *testing.T) {
		parentUUID := uuid.New()
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		handler := &CreateKategoriCommandHandler{Repo: repo}
		parentUUIDStr := parentUUID.String()
		cmd := CreateKategoriCommand{
			NamaKategori: "Child Category",
			SubKategori:  &parentUUIDStr,
			Resource:     "local",
			SID:          "sid",
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainkategori.NotFound(parentUUIDStr), err)
	})

	t.Run("Failure case - database error on get parent", func(t *testing.T) {
		parentUUID := uuid.New()
		dbErr := errors.New("db error")
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				return nil, dbErr
			},
		}

		handler := &CreateKategoriCommandHandler{Repo: repo}
		parentUUIDStr := parentUUID.String()
		cmd := CreateKategoriCommand{
			NamaKategori: "Child Category",
			SubKategori:  &parentUUIDStr,
			Resource:     "local",
			SID:          "sid",
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})

	t.Run("Failure case - invalid kategori data (validation failure)", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{}
		handler := &CreateKategoriCommandHandler{Repo: repo}
		cmd := CreateKategoriCommand{
			NamaKategori: "", // empty title, will fail domain validation
			Resource:     "invalid-owner",
			SID:          "sid",
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainkategori.InvalidOwner(), err)
	})

	t.Run("Failure case - BeginTx fails", func(t *testing.T) {
		txErr := errors.New("tx error")
		repo := &mock.MockKategoriRepository{
			BeginTxFunc: func(ctx context.Context) (*gorm.DB, error) {
				return nil, txErr
			},
		}

		handler := &CreateKategoriCommandHandler{Repo: repo}
		cmd := CreateKategoriCommand{
			NamaKategori: "Test Category",
			Resource:     "local",
			SID:          "sid",
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, txErr, err)
	})

	t.Run("Failure case - Create fails", func(t *testing.T) {
		createErr := errors.New("create failed")
		repo := &mock.MockKategoriRepository{
			CreateFunc: func(ctx context.Context, kategori *domainkategori.Kategori) error {
				return createErr
			},
		}

		handler := &CreateKategoriCommandHandler{Repo: repo}
		cmd := CreateKategoriCommand{
			NamaKategori: "Test Category",
			Resource:     "local",
			SID:          "sid",
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, createErr, err)
	})
}
