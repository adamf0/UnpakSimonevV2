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
	FullTexts       *string
	CreatedBy       *string
	CreatedByRef    *string
	DeletedAt       *time.Time
}
