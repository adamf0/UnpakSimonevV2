package application

import (
	"context"

	domainTemplateJawaban "UnpakSiamida/modules/templatejawaban/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetTemplateJawabanDefaultByUuidQueryHandler struct {
	Repo domainTemplateJawaban.ITemplateJawabanRepository
}

func (h *GetTemplateJawabanDefaultByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetTemplateJawabanDefaultByUuidQuery,
) (*domainTemplateJawaban.TemplateJawabanDefault, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainTemplateJawaban.NotFound(q.Uuid)
	}

	TemplateJawaban, err := h.Repo.GetDefaultByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainTemplateJawaban.NotFound(q.Uuid)
		}
		return nil, err
	}

	return TemplateJawaban, nil
}
