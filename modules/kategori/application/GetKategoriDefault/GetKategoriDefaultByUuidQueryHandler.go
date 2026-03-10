package application

import (
	"context"

	domainKategori "UnpakSiamida/modules/kategori/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetKategoriDefaultByUuidQueryHandler struct {
	Repo domainKategori.IKategoriRepository
}

func (h *GetKategoriDefaultByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetKategoriDefaultByUuidQuery,
) (*domainKategori.KategoriDefault, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainKategori.NotFound(q.Uuid)
	}

	Kategori, err := h.Repo.GetDefaultByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainKategori.NotFound(q.Uuid)
		}
		return nil, err
	}

	return Kategori, nil
}
