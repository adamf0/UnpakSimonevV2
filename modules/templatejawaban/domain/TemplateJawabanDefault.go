package domain

import (
	"time"

	"github.com/google/uuid"
)

type TemplateJawabanDefault struct {
	ID                     uint `json:"-"`
	UUID                   uuid.UUID
	IdTemplatePertanyaan   *uint `json:"-"`
	UUIDTemplatePertanyaan *uuid.UUID
	NamaTemplatePertanyaan *string
	Jawaban                string
	Nilai                  uint
	IsFreeText             uint
	DeletedAt              *time.Time
	CreatedAt              time.Time
	UpdatedAt              *time.Time
}
