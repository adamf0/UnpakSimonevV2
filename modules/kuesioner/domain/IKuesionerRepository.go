package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type IKuesionerRepository interface {
	GetReportSummaryOverview(ctx context.Context, rawJudul, kodeFakultas, kodeProdi, unit string) (*ReportSummaryOverview, error)
	GetReportDistribusiFakultas(ctx context.Context, rawJudul, kodeFakultas, kodeProdi, unit string) ([]ReportDistribusiFakultas, error)
	GetReportTopQuestions(ctx context.Context, rawJudul, kodeFakultas, kodeProdi, unit string) ([]ReportTopQuestion, error)
	GetReportKategoriSummary(ctx context.Context, rawJudul, kodeFakultas, kodeProdi, unit string) ([]ReportKategoriSummary, error)
	GetReportYear(ctx context.Context) ([]ReportYear, error)
	GetDashboardStats(ctx context.Context) (*DashboardStats, error)
	GetReportSummary(ctx context.Context, rawJudul, kodeFakultas, kodeProdi, unit string) (*ReportSummaryResponse, error)
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
