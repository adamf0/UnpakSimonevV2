package domain

import (
	"time"

	"github.com/google/uuid"
)

type BankSoalDefault struct {
	Id               uint `json:"-"`
	UUID             uuid.UUID
	Judul            string
	Content          *string
	Deskripsi        *string
	Semester         *string
	TanggalMulai     *time.Time
	TanggalAkhir     *time.Time
	CreatedBy        *string
	CreatedByRef     *string
	DeletedAt        *time.Time
	Status           string
	KodeFakultas     string
	KodeProdi        string
	NamaFakultas     string
	NamaProdi        string
	Role             string
	TotalPertanyaan  uint
	TotalInput       uint
	TargetPertanyaan []uuid.UUID `gorm:"-"`
	RawTargetUUIDs   string      `gorm:"column:TargetPertanyaan" json:"-"`
	UUIDKuesioner    uuid.UUID
	ListExt          []BankSoalExtDefault `gorm:"-" json:"ListExt"`
}
