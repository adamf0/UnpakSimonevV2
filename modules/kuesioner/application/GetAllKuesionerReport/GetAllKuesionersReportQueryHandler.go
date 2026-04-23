package application

import (
	domainKuesioner "UnpakSiamida/modules/kuesioner/domain"
	"context"
	"time"
)

type GetAllKuesionersReportQueryHandler struct {
	Repo domainKuesioner.IKuesionerRepository
}

func (h *GetAllKuesionersReportQueryHandler) Handle(
	ctx context.Context,
	q GetAllKuesionersReportQuery,
) ([]domainKuesioner.KuesionerResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	Kuesioners, err := h.Repo.GetAllKuesionerResult(
		ctx,
		q.JudulBankSoal,
		q.Semester,
		q.Is4Year,
	)
	if err != nil {
		return []domainKuesioner.KuesionerResult{}, err
	}

	return Kuesioners, nil
}
