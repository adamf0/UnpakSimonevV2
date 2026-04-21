package domain

import "github.com/google/uuid"

type KuesionerJawabanDefault struct {
	ID                     uint
	UUID                   *string
	IdKuesioner            uint
	UuidKuesioner          *uuid.UUID
	IdTemplatePertanyaan   uint
	UuidTemplatePertanyaan *uuid.UUID
	IdTemplateJawaban      *uint
	UuidTemplateJawaban    *uuid.UUID
	FreeText               *string
}
