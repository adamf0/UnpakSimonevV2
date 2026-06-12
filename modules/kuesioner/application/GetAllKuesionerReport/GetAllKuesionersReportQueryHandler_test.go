package application

import (
	"context"
	"errors"
	"sync"
	"testing"

	mockkuesioner "UnpakSiamida/modules/kuesioner/application/mock"
	domainkuesioner "UnpakSiamida/modules/kuesioner/domain"

	"github.com/stretchr/testify/assert"
)

func TestGetAllKuesionersReportQueryHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var mu sync.Mutex
		calledPartitions := make(map[string]bool)

		repo := &mockkuesioner.MockKuesionerRepository{
			GetAllKuesionerResultFunc: func(ctx context.Context, JudulBankSoal *string, Is4Year bool, PartitionKey string) ([]domainkuesioner.KuesionerResult, error) {
				mu.Lock()
				calledPartitions[PartitionKey] = true
				mu.Unlock()

				return []domainkuesioner.KuesionerResult{
					{
						UUID:       "uuid-" + PartitionKey,
						Pertanyaan: "Q-" + PartitionKey,
					},
				}, nil
			},
		}

		handler := &GetAllKuesionersReportQueryHandler{
			Repo: repo,
		}

		judul := "Kuesioner 2024"
		q := GetAllKuesionersReportQuery{
			JudulBankSoal: &judul,
			Is4Year:       false,
		}

		res, err := handler.Handle(context.Background(), q)
		assert.NoError(t, err)
		assert.Len(t, res, 10)

		// Assert all partitions were called
		expectedPartitions := []string{"UNIT", "01", "02", "03", "04", "05", "06", "07", "08", "09"}
		for _, pk := range expectedPartitions {
			assert.True(t, calledPartitions[pk], "expected partition %s to be called", pk)
		}
	})

	t.Run("Fail on Partition Error", func(t *testing.T) {
		expectedErr := errors.New("partition query failed")

		repo := &mockkuesioner.MockKuesionerRepository{
			GetAllKuesionerResultFunc: func(ctx context.Context, JudulBankSoal *string, Is4Year bool, PartitionKey string) ([]domainkuesioner.KuesionerResult, error) {
				if PartitionKey == "03" {
					return nil, expectedErr
				}
				return []domainkuesioner.KuesionerResult{
					{
						UUID: "uuid-" + PartitionKey,
					},
				}, nil
			},
		}

		handler := &GetAllKuesionersReportQueryHandler{
			Repo: repo,
		}

		judul := "Kuesioner 2024"
		q := GetAllKuesionersReportQuery{
			JudulBankSoal: &judul,
			Is4Year:       false,
		}

		_, err := handler.Handle(context.Background(), q)
		assert.ErrorIs(t, err, expectedErr)
	})
}
