package domain

import "github.com/google/uuid"

type KuesionerJawabanDefault struct {
	ID                     uint `json:"-"`
	UUID                   *string
	IdKuesioner            uint
	UuidKuesioner          *uuid.UUID
	IdTemplatePertanyaan   uint `json:"-"`
	UuidTemplatePertanyaan *uuid.UUID
	IdTemplateJawaban      *uint `json:"-"`
	UuidTemplateJawaban    *uuid.UUID
	FreeText               *string
}
