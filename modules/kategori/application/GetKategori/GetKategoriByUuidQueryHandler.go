package application

import (
	"context"

	domainKategori "UnpakSiamida/modules/kategori/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetKategoriByUuidQueryHandler struct {
	Repo domainKategori.IKategoriRepository
}

func (h *GetKategoriByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetKategoriByUuidQuery,
) (*domainKategori.Kategori, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainKategori.NotFound(q.Uuid)
	}

	inKategori, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainKategori.NotFound(q.Uuid)
		}
		return nil, err
	}

	return inKategori, nil
}
