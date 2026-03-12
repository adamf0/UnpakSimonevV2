package application

import (
	"context"

	domainKuesioner "UnpakSiamida/modules/kuesioner/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetKuesionerByUuidQueryHandler struct {
	Repo domainKuesioner.IKuesionerRepository
}

func (h *GetKuesionerByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetKuesionerByUuidQuery,
) (*domainKuesioner.Kuesioner, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainKuesioner.NotFound(q.Uuid)
	}

	inKuesioner, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainKuesioner.NotFound(q.Uuid)
		}
		return nil, err
	}

	return inKuesioner, nil
}
