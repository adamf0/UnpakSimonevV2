package domain

import (
	"fmt"
	"time"

	common "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"

	"github.com/google/uuid"
)

type BankSoal struct {
	common.Entity

	ID           uint       `gorm:"primaryKey;autoIncrement"`
	UUID         uuid.UUID  `gorm:"type:char(36);uniqueIndex"`
	Judul        string     `gorm:"column:judul"`
	Content      *string    `gorm:"column:content"`
	Deskripsi    *string    `gorm:"column:deskripsi"`
	Semester     *string    `gorm:"column:semester"`
	TanggalMulai *time.Time `gorm:"column:tanggal_mulai"`
	TanggalAkhir *time.Time `gorm:"column:tanggal_akhir"`
	CreatedBy    *string    `gorm:"column:createdBy"`
	CreatedByRef *string    `gorm:"column:createdByRef"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
	Status       string
}

func (BankSoal) TableName() string {
	return "bank_soalv2"
}

// === CREATE ===
func NewBankSoal(
	judul string,
	content *string,
	deskripsi *string,
	semester *string,
	tanggalmulai *string,
	tanggalakhir *string,
	createdby string, //lpm, fakultas, prodi
	createdbyref string,
) common.ResultValue[*BankSoal] {
	format := "2006-01-02 15:04:05"

	var tanggalMulai time.Time
	var tanggalAkhir time.Time
	var err error

	if createdby != "local" {
		return common.FailureValue[*BankSoal](InvalidOwner())
	}

	if tanggalmulai != nil {
		tanggalMulai, err = time.Parse(format, helper.StringValue(tanggalmulai))
		if err != nil {
			return common.FailureValue[*BankSoal](InvalidDate("tanggal mulai"))
		}
	}

	if tanggalakhir != nil {
		tanggalAkhir, err = time.Parse(format, helper.StringValue(tanggalakhir))
		if err != nil {
			return common.FailureValue[*BankSoal](InvalidDate("tanggal akhir"))
		}
	}

	if tanggalmulai != nil && tanggalakhir != nil {
		if isOverlap(tanggalMulai, tanggalAkhir) {
			return common.FailureValue[*BankSoal](InvalidDateRange())
		}
	}

	aktivitasproker := &BankSoal{
		UUID:         uuid.New(),
		Judul:        judul,
		Content:      content,
		Deskripsi:    deskripsi,
		Semester:     semester,
		TanggalMulai: &tanggalMulai,
		TanggalAkhir: &tanggalAkhir,
		Status:       "draf",
		CreatedBy:    helper.StrPtr(createdby),
		CreatedByRef: helper.StrPtr(createdbyref),
	}

	return common.SuccessValue(aktivitasproker)
}

// === UPDATE ===
func UpdateBankSoal(
	prev *BankSoal,
	uid uuid.UUID,
	judul string,
	content *string,
	deskripsi *string,
	semester *string,
	tanggalmulai *string,
	tanggalakhir *string,
	createdby string, //lpm, fakultas, prodi
	createdbyref string,
) common.ResultValue[*BankSoal] {

	if prev == nil {
		return common.FailureValue[*BankSoal](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*BankSoal](InvalidData())
	}

	format := "2006-01-02 15:04:05"

	var tanggalMulai time.Time
	var tanggalAkhir time.Time
	var err error

	if createdby != "local" {
		return common.FailureValue[*BankSoal](InvalidOwner())
	}

	if tanggalmulai != nil {
		tanggalMulai, err = time.Parse(format, helper.StringValue(tanggalmulai))
		if err != nil {
			return common.FailureValue[*BankSoal](InvalidDate("tanggal rk awal"))
		}
	}

	if tanggalakhir != nil {
		tanggalAkhir, err = time.Parse(format, helper.StringValue(tanggalakhir))
		if err != nil {
			return common.FailureValue[*BankSoal](InvalidDate("tanggal rk akhir"))
		}
	}

	if tanggalmulai != nil && tanggalakhir != nil {
		if isOverlap(tanggalMulai, tanggalAkhir) {
			return common.FailureValue[*BankSoal](InvalidDateRange())
		}
	}

	prev.Judul = judul
	prev.Content = content
	prev.Deskripsi = deskripsi
	prev.Semester = semester
	prev.TanggalMulai = &tanggalMulai
	prev.TanggalAkhir = &tanggalAkhir
	prev.CreatedBy = helper.StrPtr(createdby)
	prev.CreatedByRef = helper.StrPtr(createdbyref)

	return common.SuccessValue(prev)
}

// === Delete ===
func DeleteBankSoal(
	prev *BankSoal,
) common.ResultValue[*BankSoal] {

	if prev == nil {
		return common.FailureValue[*BankSoal](EmptyData())
	}

	now := time.Now()
	prev.DeletedAt = &now

	return common.SuccessValue(prev)
}

// === Restore ===
func RestoreBankSoal(
	prev *BankSoal,
) common.ResultValue[*BankSoal] {

	if prev == nil {
		return common.FailureValue[*BankSoal](EmptyData())
	}

	prev.DeletedAt = nil

	return common.SuccessValue(prev)
}

// === Copy ===
func CopyBankSoal(
	prev *BankSoal,
	copyCount int,
	createdby string, //lpm, fakultas, prodi
	createdbyref string,
) common.ResultValue[*BankSoal] {

	if prev == nil {
		return common.FailureValue[*BankSoal](EmptyData())
	}

	var judul string
	if copyCount == 0 {
		judul = fmt.Sprintf("salin - %s", prev.Judul)
	} else {
		judul = fmt.Sprintf("salin (%d) - %s", copyCount+1, prev.Judul)
	}

	aktivitasproker := &BankSoal{
		UUID:         uuid.New(),
		Judul:        judul,
		Content:      prev.Content,
		Deskripsi:    prev.Deskripsi,
		Semester:     prev.Semester,
		TanggalMulai: prev.TanggalMulai,
		TanggalAkhir: prev.TanggalAkhir,
		Status:       prev.Status,
		CreatedBy:    helper.StrPtr(createdby),
		CreatedByRef: helper.StrPtr(createdbyref),
	}

	return common.SuccessValue(aktivitasproker)
}

func isOverlap(start1, end1 time.Time) bool {
	return !end1.After(start1)
}
