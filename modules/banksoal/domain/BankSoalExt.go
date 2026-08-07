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

	var tanggalMulaiPtr *time.Time
	var tanggalAkhirPtr *time.Time

	if createdby != "local" {
		return common.FailureValue[*BankSoalExt](InvalidOwner())
	}

	if tanggalmulai != nil && helper.StringValue(tanggalmulai) != "" {
		t, err := parseAndAdjustDate(helper.StringValue(tanggalmulai), false)
		if err != nil {
			return common.FailureValue[*BankSoalExt](InvalidDate("tanggal awal"))
		}
		tanggalMulaiPtr = &t
	}

	if tanggalakhir != nil && helper.StringValue(tanggalakhir) != "" {
		t, err := parseAndAdjustDate(helper.StringValue(tanggalakhir), true)
		if err != nil {
			return common.FailureValue[*BankSoalExt](InvalidDate("tanggal akhir"))
		}
		tanggalAkhirPtr = &t
	}

	if tanggalMulaiPtr != nil && tanggalAkhirPtr != nil {
		if isOverlap(*tanggalMulaiPtr, *tanggalAkhirPtr) {
			return common.FailureValue[*BankSoalExt](InvalidDateRange())
		}
	}

	banksoalext := &BankSoalExt{
		UUID:         uuid.New(),
		IdBankSoal:   banksoal.ID,
		TanggalMulai: tanggalMulaiPtr,
		TanggalAkhir: tanggalAkhirPtr,
		CreatedBy:    helper.StrPtr(createdby),
		CreatedByRef: helper.StrPtr(createdbyref),
	}

	return common.SuccessValue(banksoalext)
}
