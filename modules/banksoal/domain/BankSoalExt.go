package domain

import (
	"time"

	common "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"

	"github.com/google/uuid"
)

type BankSoalExt struct {
	common.Entity

	ID           uint       `gorm:"primaryKey;autoIncrement"`
	UUID         uuid.UUID  `gorm:"type:char(36);uniqueIndex"`
	IdBankSoal   uint       `gorm:"column:id_bank_soal"`
	TanggalMulai *time.Time `gorm:"column:tanggal_mulai"`
	TanggalAkhir *time.Time `gorm:"column:tanggal_akhir"`
	CreatedBy    *string    `gorm:"column:createdBy"`
	CreatedByRef *string    `gorm:"column:createdByRef"`
}

func (BankSoalExt) TableName() string {
	return "bank_soal_extendv2"
}

func AddTimeBankSoalExt(
	banksoal *BankSoal,
	uid uuid.UUID,
	tanggalmulai *string,
	tanggalakhir *string,
	createdby string, //lpm, fakultas, prodi
	createdbyref string,
) common.ResultValue[*BankSoalExt] {

	if banksoal == nil {
		return common.FailureValue[*BankSoalExt](InvalidData())
	}

	format := "2006-01-02"

	var tanggalMulai time.Time
	var tanggalAkhir time.Time
	var err error

	if createdby != "local" {
		return common.FailureValue[*BankSoalExt](InvalidOwner())
	}

	if tanggalmulai != nil {
		tanggalMulai, err = time.Parse(format, helper.StringValue(tanggalmulai))
		if err != nil {
			return common.FailureValue[*BankSoalExt](InvalidDate("tanggal awal"))
		}
	}

	if tanggalakhir != nil {
		tanggalAkhir, err = time.Parse(format, helper.StringValue(tanggalakhir))
		if err != nil {
			return common.FailureValue[*BankSoalExt](InvalidDate("tanggal akhir"))
		}
	}

	if tanggalmulai != nil && tanggalakhir != nil {
		if isOverlap(tanggalMulai, tanggalAkhir) {
			return common.FailureValue[*BankSoalExt](InvalidDateRange())
		}
	}

	banksoalext := &BankSoalExt{
		UUID:         uuid.New(),
		IdBankSoal:   banksoal.ID,
		TanggalMulai: &tanggalMulai,
		TanggalAkhir: &tanggalAkhir,
		CreatedBy:    helper.StrPtr(createdby),
		CreatedByRef: helper.StrPtr(createdbyref),
	}

	return common.SuccessValue(banksoalext)
}
