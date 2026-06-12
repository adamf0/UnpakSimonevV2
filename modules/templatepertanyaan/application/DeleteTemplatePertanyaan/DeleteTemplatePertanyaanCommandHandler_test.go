package application_test

import (
	"context"
	"errors"
	"testing"

	"UnpakSiamida/modules/templatepertanyaan/application/DeleteTemplatePertanyaan"
	"UnpakSiamida/modules/templatepertanyaan/application/mock"
	"UnpakSiamida/modules/templatepertanyaan/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteTemplatePertanyaanCommandHandler_SoftDeleteSuccess(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	existingTp := &domain.TemplatePertanyaan{
		ID:         1,
		UUID:       uid,
		Pertanyaan: "Pertanyaan Test",
	}

	repo.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaan, error) {
		assert.Equal(t, uid, id)
		return existingTp, nil
	}

	repo.UpdateFunc = func(ctx context.Context, tp *domain.TemplatePertanyaan) error {
		assert.NotNil(t, tp.DeletedAt)
		return nil
	}

	handler := &application.DeleteTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.DeleteTemplatePertanyaanCommand{
		Uuid: uid.String(),
		Mode: "soft_delete",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, uid.String(), res)
}

func TestDeleteTemplatePertanyaanCommandHandler_HardDeleteSuccess(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	existingTp := &domain.TemplatePertanyaan{
		ID:         1,
		UUID:       uid,
		Pertanyaan: "Pertanyaan Test",
	}

	repo.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaan, error) {
		return existingTp, nil
	}

	repo.DeleteFunc = func(ctx context.Context, id uuid.UUID) error {
		assert.Equal(t, uid, id)
		return nil
	}

	handler := &application.DeleteTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.DeleteTemplatePertanyaanCommand{
		Uuid: uid.String(),
		Mode: "hard_delete",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, uid.String(), res)
}

func TestDeleteTemplatePertanyaanCommandHandler_InvalidUuid(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}

	handler := &application.DeleteTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.DeleteTemplatePertanyaanCommand{
		Uuid: "invalid-uuid",
		Mode: "soft_delete",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "uuid is invalid")
	assert.Empty(t, res)
}

func TestDeleteTemplatePertanyaanCommandHandler_NotFound(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	repo.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaan, error) {
		return nil, gorm.ErrRecordNotFound
	}

	handler := &application.DeleteTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.DeleteTemplatePertanyaanCommand{
		Uuid: uid.String(),
		Mode: "soft_delete",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Empty(t, res)
}

func TestDeleteTemplatePertanyaanCommandHandler_RepoErrorOnDelete(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	existingTp := &domain.TemplatePertanyaan{
		ID:         1,
		UUID:       uid,
		Pertanyaan: "Pertanyaan Test",
	}

	repo.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaan, error) {
		return existingTp, nil
	}

	repo.DeleteFunc = func(ctx context.Context, id uuid.UUID) error {
		return errors.New("db error")
	}

	handler := &application.DeleteTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.DeleteTemplatePertanyaanCommand{
		Uuid: uid.String(),
		Mode: "hard_delete",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	assert.Empty(t, res)
}
