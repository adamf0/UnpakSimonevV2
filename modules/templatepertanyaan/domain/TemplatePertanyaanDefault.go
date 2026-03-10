package domain

import (
	"time"

	"github.com/google/uuid"
)

type TemplatePertanyaanDefault struct {
	ID           uint
	UUID         uuid.UUID
	IdBankSoal   *uint
	UUIDBankSoal *uuid.UUID
	NamaBankSoal *string
	Pertanyaan   string
	JenisPilihan string
	Bobot        int
	IdKategori   *int
	UuidKategori *uuid.UUID
	Kategori     *string
	Required     int
	CreatedBy    *string
	CreatedByRef *string
	Fakultas     *string
	Prodi        *string
	Unit         *string
	Jenjang      *string
	DeletedAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}
