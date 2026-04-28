package application

import (
	"context"

	commondomain "UnpakSiamida/common/domain"
	domainBankSoal "UnpakSiamida/modules/banksoal/domain"
	domainTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetTemplatePertanyaanWithAnswareDefaultByBankSoalQueryHandler struct {
	Repo         domainTemplatePertanyaan.ITemplatePertanyaanRepository
	RepoBankSoal domainBankSoal.IBankSoalRepository
}

func (h *GetTemplatePertanyaanWithAnswareDefaultByBankSoalQueryHandler) Handle(
	ctx context.Context,
	q GetTemplatePertanyaanWithAnswareDefaultByBankSoalQuery,
) (commondomain.Paged[domainTemplatePertanyaan.TemplatePertanyaanWithAnswareDefault], error) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.UuidBankSoal)
	if err != nil {
		return commondomain.Paged[domainTemplatePertanyaan.TemplatePertanyaanWithAnswareDefault]{}, domainTemplatePertanyaan.InvalidBankSoal()
	}

	banksoal, err := h.RepoBankSoal.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return commondomain.Paged[domainTemplatePertanyaan.TemplatePertanyaanWithAnswareDefault]{}, domainTemplatePertanyaan.NotFoundBankSoal()
		}
		return commondomain.Paged[domainTemplatePertanyaan.TemplatePertanyaanWithAnswareDefault]{}, err
	}

	result, err := h.Repo.GetDefaultWithAnswareByBankSoal(ctx, banksoal.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return commondomain.Paged[domainTemplatePertanyaan.TemplatePertanyaanWithAnswareDefault]{}, domainTemplatePertanyaan.NotFound(q.UuidBankSoal)
		}
		return commondomain.Paged[domainTemplatePertanyaan.TemplatePertanyaanWithAnswareDefault]{}, err
	}

	return commondomain.Paged[domainTemplatePertanyaan.TemplatePertanyaanWithAnswareDefault]{
		Data:        result,
		Total:       int64(len(result)),
		CurrentPage: 1,
		TotalPages:  1,
	}, nil
}
