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

	var tanggalMulai time.Time
	var tanggalAkhir time.Time
	var err error

	if tanggalmulai != nil && helper.StringValue(tanggalmulai) != "" {
		tanggalMulai, err = parseAndAdjustDate(helper.StringValue(tanggalmulai), false)
		if err != nil {
			return common.FailureValue[*BankSoal](InvalidDate("tanggal awal"))
		}
		prev.TanggalMulai = &tanggalMulai
	}

	if tanggalakhir != nil && helper.StringValue(tanggalakhir) != "" {
		tanggalAkhir, err = parseAndAdjustDate(helper.StringValue(tanggalakhir), true)
		if err != nil {
			return common.FailureValue[*BankSoal](InvalidDate("tanggal akhir"))
		}
		prev.TanggalAkhir = &tanggalAkhir
	}

	if prev.TanggalMulai != nil && prev.TanggalAkhir != nil {
		if isOverlap(*prev.TanggalMulai, *prev.TanggalAkhir) {
			return common.FailureValue[*BankSoal](InvalidDateRange())
		}
	}

	return common.SuccessValue(prev)
}

func parseAndAdjustDate(str string, isEnd bool) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	var t time.Time
	var err error
	for _, f := range formats {
		t, err = time.Parse(f, str)
		if err == nil {
			if isEnd {
				return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location()), nil
			}
			return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()), nil
		}
	}
	return t, err
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
