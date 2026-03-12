package domain

import (
	"time"

	common "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"

	"github.com/google/uuid"
)

type TemplateJawaban struct {
	common.Entity

	ID                   uint       `gorm:"primaryKey;autoIncrement"`
	UUID                 uuid.UUID  `gorm:"type:char(36);uniqueIndex"`
	IdTemplatePertanyaan uint       `gorm:"column:id_template_pertanyaan"`
	Jawaban              string     `gorm:"column:jawaban"`
	Nilai                uint       `gorm:"column:nilai"`
	IsFreeText           uint       `gorm:"column:isFreeText"`
	CreatedBy            *string    `gorm:"column:createdBy"`
	CreatedByRef         *string    `gorm:"column:createdByRef"`
	DeletedAt            *time.Time `gorm:"column:deleted_at"`
	CreatedAt            time.Time  `gorm:"column:created_at"`
	UpdatedAt            *time.Time `gorm:"column:updated_at"`
}

func (TemplateJawaban) TableName() string {
	return "template_pilihanv2"
}

// === CREATE ===
func NewTemplateJawaban(
	id_template_pertanyaan uint,
	jawaban string,
	nilai uint,
	isFreeText uint,
	createdby string,
	createdbyref string,
) common.ResultValue[*TemplateJawaban] {

	if createdby != "local" {
		return common.FailureValue[*TemplateJawaban](InvalidOwner())
	}

	entity := &TemplateJawaban{
		UUID:                 uuid.New(),
		IdTemplatePertanyaan: id_template_pertanyaan,
		Jawaban:              jawaban,
		Nilai:                nilai,
		IsFreeText:           isFreeText,
		CreatedBy:            helper.StrPtr(createdby),
		CreatedByRef:         helper.StrPtr(createdbyref),
	}

	return common.SuccessValue(entity)
}

// === UPDATE ===
func UpdateTemplateJawaban(
	prev *TemplateJawaban,
	uid uuid.UUID,
	id_template_pertanyaan uint,
	jawaban string,
	nilai uint,
	isFreeText uint,
	createdby string,
	createdbyref string,
) common.ResultValue[*TemplateJawaban] {

	if prev == nil {
		return common.FailureValue[*TemplateJawaban](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*TemplateJawaban](InvalidData())
	}

	prev.IdTemplatePertanyaan = id_template_pertanyaan
	prev.Jawaban = jawaban
	prev.Nilai = nilai
	prev.IsFreeText = isFreeText
	prev.CreatedBy = helper.StrPtr(createdby)
	prev.CreatedByRef = helper.StrPtr(createdbyref)

	return common.SuccessValue(prev)
}

// === Delete ===
func DeleteTemplateJawaban(
	prev *TemplateJawaban,
) common.ResultValue[*TemplateJawaban] {

	if prev == nil {
		return common.FailureValue[*TemplateJawaban](EmptyData())
	}

	now := time.Now()
	prev.DeletedAt = &now

	return common.SuccessValue(prev)
}

// === Restore ===
func RestoreTemplateJawaban(
	prev *TemplateJawaban,
) common.ResultValue[*TemplateJawaban] {

	if prev == nil {
		return common.FailureValue[*TemplateJawaban](EmptyData())
	}

	prev.DeletedAt = nil

	return common.SuccessValue(prev)
}

func isOverlap(start1, end1 time.Time) bool {
	return !end1.After(start1)
}
