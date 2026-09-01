package application

import (
	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	domainBankSoal "UnpakSiamida/modules/banksoal/domain"
	"context"
	"time"

	"github.com/goforj/godump"
)

type GetAllBankSoalsQueryHandler struct {
	Repo domainBankSoal.IBankSoalRepository
}

func (h *GetAllBankSoalsQueryHandler) Handle(
	ctx context.Context,
	q GetAllBankSoalsQuery,
) (commondomain.Paged[domainBankSoal.BankSoalDefault], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	godump.Dump(q)
	if q.NIDN != nil || q.NIP != nil || q.NPM != nil {
		return commondomain.Paged[domainBankSoal.BankSoalDefault]{}, domainBankSoal.OnlyAdminFacultyStudyProgram()
	}

	BankSoals, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilter,
		helper.NullableString(q.TargetFakultas),
		helper.NullableString(q.TargetProdi),
		"",
		"",
		helper.NullableString(q.UserRole),
		q.Page,
		q.Limit,
		q.Deleted,
		false,
	)
	if err != nil {
		return commondomain.Paged[domainBankSoal.BankSoalDefault]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainBankSoal.BankSoalDefault]{
		Data:        BankSoals,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
