package application

import (
	"context"

	domainTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetTemplatePertanyaanByUuidQueryHandler struct {
	Repo domainTemplatePertanyaan.ITemplatePertanyaanRepository
}

func (h *GetTemplatePertanyaanByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetTemplatePertanyaanByUuidQuery,
) (*domainTemplatePertanyaan.TemplatePertanyaan, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainTemplatePertanyaan.NotFound(q.Uuid)
	}

	inTemplatePertanyaan, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainTemplatePertanyaan.NotFound(q.Uuid)
		}
		return nil, err
	}

	return inTemplatePertanyaan, nil
}
