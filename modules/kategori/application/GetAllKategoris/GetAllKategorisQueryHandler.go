package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainKategori "UnpakSiamida/modules/kategori/domain"
	"context"
	"time"
)

type GetAllKategorisQueryHandler struct {
	Repo domainKategori.IKategoriRepository
}

func (h *GetAllKategorisQueryHandler) Handle(
	ctx context.Context,
	q GetAllKategorisQuery,
) (commondomain.Paged[domainKategori.KategoriDefault], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	Kategoris, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
		q.Deleted,
	)
	if err != nil {
		return commondomain.Paged[domainKategori.KategoriDefault]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainKategori.KategoriDefault]{
		Data:        Kategoris,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
