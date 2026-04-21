package domain

import (
	"time"

	"github.com/google/uuid"
)

type BankSoalExtDefault struct {
	ID           uint `json:"-"`
	UUID         uuid.UUID
	IdBankSoal   uint `json:"-"`
	TanggalMulai *time.Time
	TanggalAkhir *time.Time
	CreatedBy    *string
	CreatedByRef *string
	KodeFakultas string
	KodeProdi    string
	NamaFakultas string
	NamaProdi    string
	Role         string
}
