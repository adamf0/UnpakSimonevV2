package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type IKuesionerRepository interface {
	GetReportSummaryOverview(ctx context.Context, rawJudul string) (*ReportSummaryOverview, error)
	GetReportDistribusiFakultas(ctx context.Context, rawJudul string) ([]ReportDistribusiFakultas, error)
	GetReportTopQuestions(ctx context.Context, rawJudul string) ([]ReportTopQuestion, error)
	GetReportKategoriSummary(ctx context.Context, rawJudul string) ([]ReportKategoriSummary, error)
	GetReportYear(ctx context.Context) ([]ReportYear, error)
	GetReportSummary(ctx context.Context, rawJudul string) (*ReportSummaryResponse, error)
	GetAllKuesionerResult(
		ctx context.Context,
		JudulBankSoal *string,
		// Semester *string,
		Is4Year bool,
		PartitionKey string,
	) ([]KuesionerResult, error)
	GetByUuid(ctx context.Context, uid uuid.UUID) (*Kuesioner, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*KuesionerDefault, error)
	GetAllFormFromActiveBankSoal(ctx context.Context, nidn string, nip string, npm string, banksoal []uint) ([]KuesionerDefault, error)
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
		deleted bool,
	) ([]KuesionerDefault, int64, error)
	Create(ctx context.Context, aktivitasproker *Kuesioner) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SetupUuid(ctx context.Context) error
}
