package application

import (
	commonDomain "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/kategori/application/mock"
	domainkategori "UnpakSiamida/modules/kategori/domain"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUpdateKategoriCommandHandler_Handle(t *testing.T) {
	targetUUID := uuid.New()
	targetUUIDStr := targetUUID.String()

	existingKategori := &domainkategori.Kategori{
		ID:           123,
		UUID:         targetUUID,
		NamaKategori: "Original Name",
	}

	t.Run("Success case - without parent", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				assert.Equal(t, targetUUID, uid)
				return existingKategori, nil
			},
			UpdateFunc: func(ctx context.Context, kategori *domainkategori.Kategori) error {
				assert.Equal(t, "Updated Name", kategori.NamaKategori)
				return nil
			},
			RebuildFullTextFunc: func(ctx context.Context) error {
				return nil
			},
		}

		handler := &UpdateKategoriCommandHandler{Repo: repo}
		cmd := UpdateKategoriCommand{
			Uuid:         targetUUIDStr,
			NamaKategori: "Updated Name",
			Resource:     "res",
			SID:          "sid",
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, targetUUIDStr, res)
	})

	t.Run("Success case - with valid parent", func(t *testing.T) {
		parentUUID := uuid.New()
		parentKategori := &domainkategori.Kategori{
			ID:           456,
			UUID:         parentUUID,
			NamaKategori: "Parent Category",
		}

		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				if uid == parentUUID {
					return parentKategori, nil
				}
				if uid == targetUUID {
					return existingKategori, nil
				}
				return nil, gorm.ErrRecordNotFound
			},
			UpdateFunc: func(ctx context.Context, kategori *domainkategori.Kategori) error {
				assert.Equal(t, uint(456), *kategori.SubKategori)
				return nil
			},
			RebuildFullTextFunc: func(ctx context.Context) error {
				return nil
			},
		}

		handler := &UpdateKategoriCommandHandler{Repo: repo}
		parentUUIDStr := parentUUID.String()
		cmd := UpdateKategoriCommand{
			Uuid:         targetUUIDStr,
			NamaKategori: "Updated Name",
			SubKategori:  &parentUUIDStr,
			Resource:     "res",
			SID:          "sid",
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, targetUUIDStr, res)
	})

	t.Run("Failure case - invalid UUID for category", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{}
		handler := &UpdateKategoriCommandHandler{Repo: repo}
		cmd := UpdateKategoriCommand{
			Uuid:         "invalid-uuid",
			NamaKategori: "Name",
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainkategori.InvalidUuid(), err)
	})

	t.Run("Failure case - invalid parent UUID format", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{}
		handler := &UpdateKategoriCommandHandler{Repo: repo}
		badUUID := "bad-uuid"
		cmd := UpdateKategoriCommand{
			Uuid:         targetUUIDStr,
			SubKategori:  &badUUID,
			NamaKategori: "Name",
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

		handler := &UpdateKategoriCommandHandler{Repo: repo}
		parentUUIDStr := parentUUID.String()
		cmd := UpdateKategoriCommand{
			Uuid:         targetUUIDStr,
			SubKategori:  &parentUUIDStr,
			NamaKategori: "Name",
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainkategori.NotFound(parentUUIDStr), err)
	})

	t.Run("Failure case - category not found", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		handler := &UpdateKategoriCommandHandler{Repo: repo}
		cmd := UpdateKategoriCommand{
			Uuid:         targetUUIDStr,
			NamaKategori: "Name",
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainkategori.NotFound(targetUUIDStr), err)
	})

	t.Run("Failure case - validation failure on update", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				// To force a validation failure on update, we trigger empty data error by passing nil as existing kategori
				return nil, nil
			},
		}

		handler := &UpdateKategoriCommandHandler{Repo: repo}
		cmd := UpdateKategoriCommand{
			Uuid:         targetUUIDStr,
			NamaKategori: "Name",
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainkategori.EmptyData(), err)
	})

	t.Run("Failure case - BeginTx fails", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
				return existingKategori, nil
			},
			BeginTxFunc: func(ctx context.Context) (*gorm.DB, error) {
				return nil, errors.New("begin tx error")
			},
		}

		handler := &UpdateKategoriCommandHandler{Repo: repo}
		cmd := UpdateKategoriCommand{
			Uuid:         targetUUIDStr,
			NamaKategori: "Updated Name",
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
	})
}

func TestUpdateKategoriOrderCommandHandler_Handle(t *testing.T) {
	uuid1 := uuid.New()
	uuid2 := uuid.New()

	kategoriList := []domainkategori.KategoriDefault{
		{
			ID:   1,
			UUID: uuid1,
		},
		{
			ID:   2,
			UUID: uuid2,
		},
	}

	t.Run("Success case", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{
			GetAllFunc: func(ctx context.Context, search string, filters []commonDomain.SearchFilter, page, limit *int, deleted bool) ([]domainkategori.KategoriDefault, int64, error) {
				return kategoriList, int64(len(kategoriList)), nil
			},
			UpdateParentBatchFunc: func(ctx context.Context, rows []domainkategori.UpdateRow) error {
				assert.Len(t, rows, 2)
				assert.Equal(t, uint(1), rows[0].ID)
				assert.Equal(t, uint(2), rows[1].ID)
				assert.Nil(t, rows[0].SubKategori)
				assert.Equal(t, uint(1), *rows[1].SubKategori)
				return nil
			},
		}

		handler := &UpdateKategoriOrderCommandHandler{Repo: repo}
		payloadJSON := fmt.Sprintf(`[
			{"uuid": "%s", "full_text": "First Category"},
			{"uuid": "%s", "uuidSub": "%s", "full_text": "Second Category"}
		]`, uuid1, uuid2, uuid1)

		res, err := handler.Handle(context.Background(), UpdateKategoriOrderCommand{Payload: payloadJSON})
		assert.NoError(t, err)
		assert.Equal(t, "ok", res)
	})

	t.Run("Failure case - invalid JSON", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{}
		handler := &UpdateKategoriOrderCommandHandler{Repo: repo}
		res, err := handler.Handle(context.Background(), UpdateKategoriOrderCommand{Payload: "invalid-json"})
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("Failure case - empty fulltext", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{}
		handler := &UpdateKategoriOrderCommandHandler{Repo: repo}
		payloadJSON := fmt.Sprintf(`[
			{"uuid": "%s", "full_text": ""}
		]`, uuid1)

		_, err := handler.Handle(context.Background(), UpdateKategoriOrderCommand{Payload: payloadJSON})
		assert.Error(t, err)
	})

	t.Run("Failure case - invalid UUID in payload", func(t *testing.T) {
		repo := &mock.MockKategoriRepository{}
		handler := &UpdateKategoriOrderCommandHandler{Repo: repo}
		payloadJSON := `[
			{"uuid": "not-a-uuid", "full_text": "Some text"}
		]`

		_, err := handler.Handle(context.Background(), UpdateKategoriOrderCommand{Payload: payloadJSON})
		assert.Error(t, err)
	})

	t.Run("Failure case - db error on GetAll", func(t *testing.T) {
		dbErr := errors.New("db error")
		repo := &mock.MockKategoriRepository{
			GetAllFunc: func(ctx context.Context, search string, filters []commonDomain.SearchFilter, page, limit *int, deleted bool) ([]domainkategori.KategoriDefault, int64, error) {
				return nil, 0, dbErr
			},
		}

		handler := &UpdateKategoriOrderCommandHandler{Repo: repo}
		payloadJSON := fmt.Sprintf(`[
			{"uuid": "%s", "full_text": "First"}
		]`, uuid1)

		_, err := handler.Handle(context.Background(), UpdateKategoriOrderCommand{Payload: payloadJSON})
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})
}
