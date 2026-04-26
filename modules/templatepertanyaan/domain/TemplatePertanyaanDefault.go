package domain

import (
	"time"

	"github.com/google/uuid"
)

type TemplatePertanyaanDefault struct {
	ID           uint `json:"-"`
	UUID         uuid.UUID
	IdBankSoal   *uint `json:"-"`
	UUIDBankSoal *uuid.UUID
	NamaBankSoal *string
	Pertanyaan   string
	JenisPilihan string
	Bobot        int
	IdKategori   *int `json:"-"`
	UuidKategori *uuid.UUID
	Kategori     *string
	FullPath     *string
	Required     int
	Status       string
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
