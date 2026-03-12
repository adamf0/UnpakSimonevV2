package domain

import (
	"time"

	"github.com/google/uuid"
)

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
