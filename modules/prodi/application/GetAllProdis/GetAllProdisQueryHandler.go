package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainProdi "UnpakSiamida/modules/prodi/domain"
	"context"
	"time"
)

type GetAllProdisQueryHandler struct {
	Repo domainProdi.IProdiRepository
}

func (h *GetAllProdisQueryHandler) Handle(
	ctx context.Context,
	q GetAllProdisQuery,
) (commondomain.Paged[domainProdi.ProdiDefault], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	Prodis, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
	)
	if err != nil {
		return commondomain.Paged[domainProdi.ProdiDefault]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainProdi.ProdiDefault]{
		Data:        Prodis,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
