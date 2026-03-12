package application

import (
	"context"

	domainKuesioner "UnpakSiamida/modules/kuesioner/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetKuesionerDefaultByUuidQueryHandler struct {
	Repo domainKuesioner.IKuesionerRepository
}

func (h *GetKuesionerDefaultByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetKuesionerDefaultByUuidQuery,
) (*domainKuesioner.KuesionerDefault, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainKuesioner.NotFound(q.Uuid)
	}

	Kuesioner, err := h.Repo.GetDefaultByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainKuesioner.NotFound(q.Uuid)
		}
		return nil, err
	}

	return Kuesioner, nil
}
