package domain

import (
	"time"

	"github.com/google/uuid"
)

type BankSoalDefault struct {
	Id           uint
	UUID         uuid.UUID
	Judul        string
	Content      *string
	Deskripsi    *string
	Semester     *string
	TanggalMulai *time.Time
	TanggalAkhir *time.Time
	CreatedBy    *string
	DeletedAt    *time.Time
	Status       string
}
