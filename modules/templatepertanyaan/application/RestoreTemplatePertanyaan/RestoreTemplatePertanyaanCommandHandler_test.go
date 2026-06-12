package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"UnpakSiamida/modules/templatepertanyaan/application/RestoreTemplatePertanyaan"
	"UnpakSiamida/modules/templatepertanyaan/application/mock"
	"UnpakSiamida/modules/templatepertanyaan/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRestoreTemplatePertanyaanCommandHandler_Success(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()
	deletedAt := time.Now()

	existingTp := &domain.TemplatePertanyaan{
		ID:        1,
		UUID:      uid,
		DeletedAt: &deletedAt,
	}

	repo.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaan, error) {
		assert.Equal(t, uid, id)
		return existingTp, nil
	}

	repo.UpdateFunc = func(ctx context.Context, tp *domain.TemplatePertanyaan) error {
		assert.Nil(t, tp.DeletedAt)
		return nil
	}

	handler := &application.RestoreTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.RestoreTemplatePertanyaanCommand{
		Uuid: uid.String(),
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, uid.String(), res)
}

func TestRestoreTemplatePertanyaanCommandHandler_InvalidUuid(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}

	handler := &application.RestoreTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.RestoreTemplatePertanyaanCommand{
		Uuid: "invalid-uuid",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "uuid is invalid")
	assert.Empty(t, res)
}

func TestRestoreTemplatePertanyaanCommandHandler_NotFound(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	repo.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaan, error) {
		return nil, gorm.ErrRecordNotFound
	}

	handler := &application.RestoreTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.RestoreTemplatePertanyaanCommand{
		Uuid: uid.String(),
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Empty(t, res)
}

func TestRestoreTemplatePertanyaanCommandHandler_UpdateFail(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()
	deletedAt := time.Now()

	existingTp := &domain.TemplatePertanyaan{
		ID:        1,
		UUID:      uid,
		DeletedAt: &deletedAt,
	}

	repo.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaan, error) {
		return existingTp, nil
	}

	repo.UpdateFunc = func(ctx context.Context, tp *domain.TemplatePertanyaan) error {
		return errors.New("db error")
	}

	handler := &application.RestoreTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.RestoreTemplatePertanyaanCommand{
		Uuid: uid.String(),
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	assert.Empty(t, res)
}
