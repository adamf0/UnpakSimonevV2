package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainKuesioner "UnpakSiamida/modules/kuesioner/domain"
	"context"
	"time"
)

type GetAllKuesionersQueryHandler struct {
	Repo domainKuesioner.IKuesionerRepository
}

func (h *GetAllKuesionersQueryHandler) Handle(
	ctx context.Context,
	q GetAllKuesionersQuery,
) (commondomain.Paged[domainKuesioner.KuesionerDefault], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	Kuesioners, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
		q.Deleted,
	)
	if err != nil {
		return commondomain.Paged[domainKuesioner.KuesionerDefault]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainKuesioner.KuesionerDefault]{
		Data:        Kuesioners,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
