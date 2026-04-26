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

	ID           uint       `gorm:"primaryKey;autoIncrement" json:"-"`
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
	createdby string, //lpm, fakultas, prodi
	createdbyref string,
) common.ResultValue[*BankSoal] {
	if createdby != "local" {
		return common.FailureValue[*BankSoal](InvalidOwner())
	}

	aktivitasproker := &BankSoal{
		UUID:         uuid.New(),
		Judul:        judul,
		Content:      content,
		Deskripsi:    deskripsi,
		Semester:     semester,
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
	createdby string, //lpm, fakultas, prodi
	createdbyref string,
) common.ResultValue[*BankSoal] {

	if prev == nil {
		return common.FailureValue[*BankSoal](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*BankSoal](InvalidData())
	}

	if createdby != "local" {
		return common.FailureValue[*BankSoal](InvalidOwner())
	}

	prev.Judul = judul
	prev.Content = content
	prev.Deskripsi = deskripsi
	prev.Semester = semester
	prev.CreatedBy = helper.StrPtr(createdby)
	prev.CreatedByRef = helper.StrPtr(createdbyref)

	return common.SuccessValue(prev)
}

func UpdateTimeBankSoal(
	prev *BankSoal,
	uid uuid.UUID,
	tanggalmulai *string,
	tanggalakhir *string,
) common.ResultValue[*BankSoal] {
	if prev == nil {
		return common.FailureValue[*BankSoal](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*BankSoal](InvalidData())
	}

	format := "2006-01-02"

	var tanggalMulai time.Time
	var tanggalAkhir time.Time
	var err error

	if tanggalmulai != nil {
		tanggalMulai, err = time.Parse(format, helper.StringValue(tanggalmulai))
		if err != nil {
			return common.FailureValue[*BankSoal](InvalidDate("tanggal awal"))
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

	prev.TanggalMulai = &tanggalMulai
	prev.TanggalAkhir = &tanggalAkhir

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

// === Reset Time ===
func ResetTimeBankSoal(
	prev *BankSoal,
) common.ResultValue[*BankSoal] {

	if prev == nil {
		return common.FailureValue[*BankSoal](EmptyData())
	}

	prev.TanggalMulai = nil
	prev.TanggalAkhir = nil

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

func ChangeStatus(
	prev *BankSoal,
	status string,
) common.ResultValue[*BankSoal] {

	if prev == nil {
		return common.FailureValue[*BankSoal](EmptyData())
	}

	validStatuses := map[string]bool{
		"draf":   true,
		"active": true,
	}

	if !validStatuses[status] {
		return common.FailureValue[*BankSoal](InvalidStatus())
	}

	prev.Status = status

	return common.SuccessValue(prev)
}

func isOverlap(start1, end1 time.Time) bool {
	return !end1.After(start1)
}
