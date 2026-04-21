package application

import (
	"UnpakSiamida/common/helper"
	domainBankSoal "UnpakSiamida/modules/banksoal/domain"
	domainKuesioner "UnpakSiamida/modules/kuesioner/domain"
	"context"
	"time"

	"github.com/google/uuid"
)

type ActiveKuesionerSingleQueryHandler struct {
	Repo         domainKuesioner.IKuesionerRepository
	RepoJawaban  domainKuesioner.IKuesionerJawabanRepository
	RepoBankSoal domainBankSoal.IBankSoalRepository
}

func (h *ActiveKuesionerSingleQueryHandler) Handle(
	ctx context.Context,
	q ActiveKuesionerSingleQuery,
) (*domainBankSoal.BankSoalDefault, error) {

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if q.NIDN == nil && q.NIP == nil && q.NPM == nil {
		return nil, domainBankSoal.OnlyStudentLecturerStaff()
	}

	uid, err := uuid.Parse(q.UUID)
	if err != nil {
		return nil, domainBankSoal.InvalidUuid()
	}

	kuesionerActive, err := h.Repo.GetDefaultByUuid(ctx, uid)
	if err != nil {
		return nil, err
	}

	if helper.NullableString(q.NIDN) != helper.NullableString(kuesionerActive.NIDN) ||
		helper.NullableString(q.NIP) != helper.NullableString(kuesionerActive.NIP) ||
		helper.NullableString(q.NPM) != helper.NullableString(kuesionerActive.NPM) {
		return nil, domainBankSoal.OnlyStudentLecturerStaff()
	}

	bankUUID, err := uuid.Parse(kuesionerActive.UUID.String())
	if err != nil {
		return nil, err
	}

	bankSoal, err := h.RepoBankSoal.GetDefaultByKuesioner(ctx, bankUUID)
	if err != nil {
		return nil, err
	}

	// inject uuid kuesioner
	bankSoal.UUIDKuesioner = kuesionerActive.UUID

	// =========================
	// TOTAL INPUT
	// =========================
	totalMap, err := h.RepoJawaban.GetTotalInputByKuesionerIDs(
		ctx,
		[]uint{kuesionerActive.Id},
	)
	if err != nil {
		return nil, err
	}

	if total, ok := totalMap[kuesionerActive.UUID.String()]; ok {
		bankSoal.TotalInput = total
	}

	return bankSoal, nil
}
