package domain

import (
	"fmt"
	"time"

	common "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"

	"github.com/google/uuid"
)

type TemplatePertanyaan struct {
	common.Entity

	ID           uint       `gorm:"primaryKey;autoIncrement"`
	UUID         uuid.UUID  `gorm:"type:char(36);uniqueIndex"`
	IdBankSoal   uint       `gorm:"column:id_bank_soal"`
	Pertanyaan   string     `gorm:"column:pertanyaan"`
	JenisPilihan string     `gorm:"column:jenis_pilihan"`
	Bobot        uint       `gorm:"column:bobot"`
	IdKategori   *uint      `gorm:"column:id_kategori"`
	Required     int        `gorm:"column:required"`
	CreatedBy    *string    `gorm:"column:createdBy"`
	CreatedByRef *string    `gorm:"column:createdByRef"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
}

func (TemplatePertanyaan) TableName() string {
	return "template_pertanyaanv2"
}

// === CREATE ===
func NewTemplatePertanyaan(
	id_bank_soal uint,
	pertanyaan string,
	jenisPilihan string,
	bobot uint,
	idKategori *uint,
	required int,
	createdby string,
	createdbyref string,
) common.ResultValue[*TemplatePertanyaan] {

	if createdby != "local" {
		return common.FailureValue[*TemplatePertanyaan](InvalidOwner())
	}

	if bobot <= 0 {
		bobot = 1
	}

	entity := &TemplatePertanyaan{
		UUID:         uuid.New(),
		IdBankSoal:   id_bank_soal,
		Pertanyaan:   pertanyaan,
		JenisPilihan: jenisPilihan,
		Bobot:        bobot,
		IdKategori:   idKategori,
		Required:     required,
		CreatedBy:    helper.StrPtr(createdby),
		CreatedByRef: helper.StrPtr(createdbyref),
	}

	return common.SuccessValue(entity)
}

// === UPDATE ===
func UpdateTemplatePertanyaan(
	prev *TemplatePertanyaan,
	uid uuid.UUID,
	id_bank_soal uint,
	pertanyaan string,
	jenisPilihan string,
	bobot uint,
	idKategori *uint,
	required int,
	createdby string,
	createdbyref string,
) common.ResultValue[*TemplatePertanyaan] {

	if prev == nil {
		return common.FailureValue[*TemplatePertanyaan](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*TemplatePertanyaan](InvalidData())
	}

	prev.IdBankSoal = id_bank_soal
	prev.Pertanyaan = pertanyaan
	prev.JenisPilihan = jenisPilihan
	prev.Bobot = bobot
	prev.IdKategori = idKategori
	prev.Required = required
	prev.CreatedBy = helper.StrPtr(createdby)
	prev.CreatedByRef = helper.StrPtr(createdbyref)

	return common.SuccessValue(prev)
}

// === Delete ===
func DeleteTemplatePertanyaan(
	prev *TemplatePertanyaan,
) common.ResultValue[*TemplatePertanyaan] {

	if prev == nil {
		return common.FailureValue[*TemplatePertanyaan](EmptyData())
	}

	now := time.Now()
	prev.DeletedAt = &now

	return common.SuccessValue(prev)
}

// === Restore ===
func RestoreTemplatePertanyaan(
	prev *TemplatePertanyaan,
) common.ResultValue[*TemplatePertanyaan] {

	if prev == nil {
		return common.FailureValue[*TemplatePertanyaan](EmptyData())
	}

	prev.DeletedAt = nil

	return common.SuccessValue(prev)
}

// === Copy ===
func CopyTemplatePertanyaan(
	prev *TemplatePertanyaan,
	copyCount int,
	createdby string, //lpm, fakultas, prodi
	createdbyref string,
) common.ResultValue[*TemplatePertanyaan] {

	if prev == nil {
		return common.FailureValue[*TemplatePertanyaan](EmptyData())
	}

	var pertanyaan string
	if copyCount == 0 {
		pertanyaan = fmt.Sprintf("salin - %s", prev.Pertanyaan)
	} else {
		pertanyaan = fmt.Sprintf("salin (%d) - %s", copyCount+1, prev.Pertanyaan)
	}

	aktivitasproker := &TemplatePertanyaan{
		UUID:         uuid.New(),
		Pertanyaan:   pertanyaan,
		JenisPilihan: prev.JenisPilihan,
		Bobot:        prev.Bobot,
		IdKategori:   prev.IdKategori,
		Required:     prev.Required,
		CreatedBy:    helper.StrPtr(createdby),
		CreatedByRef: helper.StrPtr(createdbyref),
	}

	return common.SuccessValue(aktivitasproker)
}

func isOverlap(start1, end1 time.Time) bool {
	return !end1.After(start1)
}
