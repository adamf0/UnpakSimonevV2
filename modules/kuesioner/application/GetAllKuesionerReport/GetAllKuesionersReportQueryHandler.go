package application

import (
	domainKuesioner "UnpakSiamida/modules/kuesioner/domain"
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

type partitionResult struct {
	data []domainKuesioner.KuesionerResult
	err  error
}

type GetAllKuesionersReportQueryHandler struct {
	Repo domainKuesioner.IKuesionerRepository
}

func (h *GetAllKuesionersReportQueryHandler) Handle(
	ctx context.Context,
	q GetAllKuesionersReportQuery,
) ([]domainKuesioner.KuesionerResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	partitionKeys := []string{
		"UNIT",
		"01",
		"02",
		"03",
		"04",
		"05",
		"06",
		"07",
		"08",
		"09",
	}

	g, ctxg := errgroup.WithContext(ctx)

	ch := make(chan partitionResult, len(partitionKeys))

	for _, pk := range partitionKeys {
		g.Go(func() error {
			rows, err := h.Repo.GetAllKuesionerResult(
				ctxg,
				q.JudulBankSoal,
				// q.Semester,
				q.Is4Year,
				pk,
			)

			ch <- partitionResult{
				data: rows,
				err:  err,
			}

			return err
		})
	}

	go func() {
		g.Wait()
		close(ch)
	}()

	var result []domainKuesioner.KuesionerResult

	for res := range ch {
		if res.err != nil {
			return nil, res.err
		}
		result = append(result, res.data...)
	}

	return result, nil
}
