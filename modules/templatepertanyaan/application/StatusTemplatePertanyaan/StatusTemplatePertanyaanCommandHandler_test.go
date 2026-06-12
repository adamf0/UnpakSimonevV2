package application_test

import (
	"context"
	"errors"
	"testing"

	"UnpakSiamida/modules/templatepertanyaan/application/StatusTemplatePertanyaan"
	"UnpakSiamida/modules/templatepertanyaan/application/mock"
	"UnpakSiamida/modules/templatepertanyaan/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestStatusTemplatePertanyaanCommandHandler_Success(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	existingTp := &domain.TemplatePertanyaan{
		ID:     1,
		UUID:   uid,
		Status: "draf",
	}

	repo.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaan, error) {
		assert.Equal(t, uid, id)
		return existingTp, nil
	}

	repo.UpdateFunc = func(ctx context.Context, tp *domain.TemplatePertanyaan) error {
		assert.Equal(t, "active", tp.Status)
		return nil
	}

	handler := &application.StatusTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.StatusTemplatePertanyaanCommand{
		Uuid:   uid.String(),
		Status: "active",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, uid.String(), res)
}

func TestStatusTemplatePertanyaanCommandHandler_InvalidUuid(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}

	handler := &application.StatusTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.StatusTemplatePertanyaanCommand{
		Uuid:   "invalid-uuid",
		Status: "active",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "uuid is invalid")
	assert.Empty(t, res)
}

func TestStatusTemplatePertanyaanCommandHandler_NotFound(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	repo.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaan, error) {
		return nil, gorm.ErrRecordNotFound
	}

	handler := &application.StatusTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.StatusTemplatePertanyaanCommand{
		Uuid:   uid.String(),
		Status: "active",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Empty(t, res)
}

func TestStatusTemplatePertanyaanCommandHandler_InvalidStatus(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	existingTp := &domain.TemplatePertanyaan{
		ID:     1,
		UUID:   uid,
		Status: "draf",
	}

	repo.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaan, error) {
		return existingTp, nil
	}

	handler := &application.StatusTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.StatusTemplatePertanyaanCommand{
		Uuid:   uid.String(),
		Status: "invalid-status",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "status is invalid")
	assert.Empty(t, res)
}

func TestStatusTemplatePertanyaanCommandHandler_UpdateFail(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	existingTp := &domain.TemplatePertanyaan{
		ID:     1,
		UUID:   uid,
		Status: "draf",
	}

	repo.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaan, error) {
		return existingTp, nil
	}

	repo.UpdateFunc = func(ctx context.Context, tp *domain.TemplatePertanyaan) error {
		return errors.New("db error")
	}

	handler := &application.StatusTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.StatusTemplatePertanyaanCommand{
		Uuid:   uid.String(),
		Status: "active",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	assert.Empty(t, res)
}
