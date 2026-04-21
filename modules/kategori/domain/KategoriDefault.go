package domain

import (
	"time"

	"github.com/google/uuid"
)

type KategoriDefault struct {
	ID              uint
	UUID            uuid.UUID
	NamaKategori    string
	IdSubKategori   *int
	UuidSubKategori *uuid.UUID
	NamaSubKategori string
	FullTexts       *string
	CreatedBy       *string
	CreatedByRef    *string
	DeletedAt       *time.Time
	KodeFakultas    string
	KodeProdi       string
	NamaFakultas    string
	NamaProdi       string
	Role            string
}
