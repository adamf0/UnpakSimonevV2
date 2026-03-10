package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"
	"context"
	"time"
)

type GetAllTemplatePertanyaansQueryHandler struct {
	Repo domainTemplatePertanyaan.ITemplatePertanyaanRepository
}

func (h *GetAllTemplatePertanyaansQueryHandler) Handle(
	ctx context.Context,
	q GetAllTemplatePertanyaansQuery,
) (commondomain.Paged[domainTemplatePertanyaan.TemplatePertanyaanDefault], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	TemplatePertanyaans, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
		q.Deleted,
	)
	if err != nil {
		return commondomain.Paged[domainTemplatePertanyaan.TemplatePertanyaanDefault]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainTemplatePertanyaan.TemplatePertanyaanDefault]{
		Data:        TemplatePertanyaans,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
