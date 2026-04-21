package application

import (
	commonDomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	domainBankSoal "UnpakSiamida/modules/banksoal/domain"
	domainKuesioner "UnpakSiamida/modules/kuesioner/domain"
	"context"
	"runtime"
	"sync"
	"time"
)

type ActiveKuesionerQueryHandler struct {
	Repo         domainKuesioner.IKuesionerRepository
	RepoJawaban  domainKuesioner.IKuesionerJawabanRepository
	RepoBankSoal domainBankSoal.IBankSoalRepository
}

func (h *ActiveKuesionerQueryHandler) Handle(
	ctx context.Context,
	q ActiveKuesionerQuery,
) ([]domainBankSoal.BankSoalDefault, error) {

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if q.NIDN == nil && q.NIP == nil && q.NPM == nil {
		return []domainBankSoal.BankSoalDefault{}, domainBankSoal.OnlyStudentLecturerStaff()
	}

	kuesionerActive, total, err := h.RepoBankSoal.GetAll(
		ctx,
		"",
		[]commonDomain.SearchFilter{},
		helper.NullableString(q.Fakultas),
		helper.NullableString(q.Prodi),
		helper.NullableString(q.Unit),
		"active",
		nil,
		nil,
		false,
		true,
	)
	if err != nil || total == 0 {
		return []domainBankSoal.BankSoalDefault{}, err
	}

	// =========================
	// Collect BankSoal IDs
	// =========================
	seen := make(map[uint]struct{})
	bankSoalKeys := make([]uint, 0, len(kuesionerActive))

	for _, k := range kuesionerActive {
		if _, ok := seen[k.Id]; !ok {
			seen[k.Id] = struct{}{}
			bankSoalKeys = append(bankSoalKeys, k.Id)
		}
	}

	// =========================
	// Get Kuesioner Forms
	// =========================
	kuesionerForms, err := h.Repo.GetAllFormFromActiveBankSoal(
		ctx,
		helper.NullableString(q.NIDN),
		helper.NullableString(q.NIP),
		helper.NullableString(q.NPM),
		bankSoalKeys,
	)
	if err != nil {
		return []domainBankSoal.BankSoalDefault{}, err
	}

	formMap := make(map[uint]domainKuesioner.KuesionerDefault)
	uuidList := make([]string, 0, len(kuesionerForms))
	idList := make([]uint, 0, len(kuesionerForms))

	for _, f := range kuesionerForms {
		idBankSoal, err := helper.ParseUint(f.IdBankSoal)
		if err != nil {
			continue
		}
		formMap[idBankSoal] = f
		uuidList = append(uuidList, f.UUID.String())
		idList = append(idList, f.Id)
	}

	// =========================
	// Get Total Input (bulk)
	// =========================
	totalInputMap := make(map[string]uint)

	if len(idList) > 0 {
		totalInputMap, err = h.RepoJawaban.GetTotalInputByKuesionerIDs(ctx, idList)
		if err != nil {
			return []domainBankSoal.BankSoalDefault{}, err
		}
	}

	// =========================
	// PARALLEL MERGE
	// =========================
	type result struct {
		index int
		data  domainBankSoal.BankSoalDefault
	}

	worker := runtime.NumCPU()
	jobs := make(chan int, len(kuesionerActive))
	results := make(chan result, len(kuesionerActive))

	var wg sync.WaitGroup

	for w := 0; w < worker; w++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for i := range jobs {
				item := kuesionerActive[i]

				if f, ok := formMap[item.Id]; ok {
					item.UUIDKuesioner = f.UUID
					if total, ok := totalInputMap[f.UUID.String()]; ok {
						item.TotalInput = total
					}
				}

				results <- result{
					index: i,
					data:  item,
				}
			}
		}()
	}

	for i := range kuesionerActive {
		jobs <- i
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	final := make([]domainBankSoal.BankSoalDefault, len(kuesionerActive))

	for res := range results {
		final[res.index] = res.data
	}

	return final, nil
}
