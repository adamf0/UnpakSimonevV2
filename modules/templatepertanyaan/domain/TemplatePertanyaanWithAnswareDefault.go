package domain

import (
	"time"

	"github.com/google/uuid"
)

type TemplatePertanyaanWithAnswareDefault struct {
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
	ListJawaban  []TemplateJawabanDefault `gorm:"-"`
}

type TemplateJawabanDefault struct {
	ID                     uint
	UUID                   uuid.UUID
	IdTemplatePertanyaan   *uint
	UUIDTemplatePertanyaan *uuid.UUID
	NamaTemplatePertanyaan *string
	Jawaban                string
	Nilai                  uint
	IsFreeText             uint
	DeletedAt              *time.Time
	CreatedAt              time.Time
	UpdatedAt              *time.Time
}
