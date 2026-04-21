package application

import (
	"context"

	domainTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetTemplatePertanyaanWithAnswareDefaultByUuidQueryHandler struct {
	Repo domainTemplatePertanyaan.ITemplatePertanyaanRepository
}

func (h *GetTemplatePertanyaanWithAnswareDefaultByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetTemplatePertanyaanWithAnswareDefaultByUuidQuery,
) (*domainTemplatePertanyaan.TemplatePertanyaanWithAnswareDefault, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainTemplatePertanyaan.NotFound(q.Uuid)
	}

	TemplatePertanyaan, err := h.Repo.GetDefaultWithAnswareByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainTemplatePertanyaan.NotFound(q.Uuid)
		}
		return nil, err
	}

	return TemplatePertanyaan, nil
}
