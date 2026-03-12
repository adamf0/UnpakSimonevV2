package application

import (
	"context"

	domainTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetTemplatePertanyaanDefaultByUuidQueryHandler struct {
	Repo domainTemplatePertanyaan.ITemplatePertanyaanRepository
}

func (h *GetTemplatePertanyaanDefaultByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetTemplatePertanyaanDefaultByUuidQuery,
) (*domainTemplatePertanyaan.TemplatePertanyaanDefault, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainTemplatePertanyaan.NotFound(q.Uuid)
	}

	TemplatePertanyaan, err := h.Repo.GetDefaultByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainTemplatePertanyaan.NotFound(q.Uuid)
		}
		return nil, err
	}

	return TemplatePertanyaan, nil
}
