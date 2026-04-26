package domain

import (
	"time"

	common "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"

	"github.com/google/uuid"
)

type KuesionerJawaban struct {
	ID                   uint       `gorm:"primaryKey;autoIncrement" json:"-"`
	UUID                 *string    `gorm:"type:varchar(36)"`
	IdKuesioner          uint       `gorm:"column:id_kuesioner"`
	IdTemplatePertanyaan uint       `gorm:"column:id_template_pertanyaan" json:"-"`
	IdTemplateJawaban    *uint      `gorm:"column:id_template_jawaban" json:"-"`
	FreeText             *string    `gorm:"column:freeText"`
	CreatedBy            *string    `gorm:"column:createdBy"`
	CreatedByRef         *string    `gorm:"column:createdByRef"`
	CreatedAt            time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt            *time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (KuesionerJawaban) TableName() string {
	return "kuesioner_jawabanv2"
}

// === CREATE ===
func NewKuesionerJawaban(
	idKuesioner uint,
	idTemplatePertanyaan uint,
	idJawaban *uint,
	freeText *string,
	createdBy *string,
	createdByRef *string,
) common.ResultValue[*KuesionerJawaban] {

	uuidStr := uuid.New().String()

	kuesioner := &KuesionerJawaban{
		UUID:                 helper.StrPtr(uuidStr),
		IdKuesioner:          idKuesioner,
		IdTemplatePertanyaan: idTemplatePertanyaan,
		IdTemplateJawaban:    idJawaban,
		FreeText:             freeText,
		CreatedBy:            createdBy,
		CreatedByRef:         createdByRef,
	}

	return common.SuccessValue(kuesioner)
}
