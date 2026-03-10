package domain

import (
	"fmt"
	"time"

	common "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"

	"github.com/google/uuid"
)

type Kategori struct {
	common.Entity

	ID           uint       `gorm:"primaryKey;autoIncrement"`
	UUID         uuid.UUID  `gorm:"type:char(36);uniqueIndex"`
	NamaKategori string     `gorm:"column:nama_kategori"`
	SubKategori  *uint      `gorm:"column:sub_kategori"`
	FullTexts    *string    `gorm:"column:full_text"`
	CreatedBy    *string    `gorm:"column:createdBy"`
	CreatedByRef *string    `gorm:"column:createdByRef"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
	CreatedAt    *time.Time `gorm:"column:created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
}

func (Kategori) TableName() string {
	return "kategori"
}

// === CREATE ===
func NewKategori(
	namaKategori string,
	subKategori *uint,
	createdby string,
	createdbyref string,
) common.ResultValue[*Kategori] {

	if createdby != "local" {
		return common.FailureValue[*Kategori](InvalidOwner())
	}

	entity := &Kategori{
		UUID:         uuid.New(),
		NamaKategori: namaKategori,
		SubKategori:  subKategori,
		CreatedBy:    helper.StrPtr(createdby),
		CreatedByRef: helper.StrPtr(createdbyref),
	}

	return common.SuccessValue(entity)
}

// === UPDATE ===
func UpdateKategori(
	prev *Kategori,
	uid uuid.UUID,
	namaKategori string,
	subKategori *uint,
	createdby string,
	createdbyref string,
) common.ResultValue[*Kategori] {

	if prev == nil {
		return common.FailureValue[*Kategori](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*Kategori](InvalidData())
	}

	prev.NamaKategori = namaKategori
	prev.SubKategori = subKategori
	prev.CreatedBy = helper.StrPtr(createdby)
	prev.CreatedByRef = helper.StrPtr(createdbyref)

	return common.SuccessValue(prev)
}

// === Move ===
func MoveKategori(
	prev *Kategori,
	parentID *uint,
) common.ResultValue[*Kategori] {

	if prev == nil {
		return common.FailureValue[*Kategori](EmptyData())
	}

	if parentID != nil && *parentID == prev.ID {
		return common.FailureValue[*Kategori](InvalidHierarchy())
	}

	prev.SubKategori = parentID

	return common.SuccessValue(prev)
}

// === Delete ===
func DeleteKategori(
	prev *Kategori,
) common.ResultValue[*Kategori] {

	if prev == nil {
		return common.FailureValue[*Kategori](EmptyData())
	}

	now := time.Now()
	prev.DeletedAt = &now

	return common.SuccessValue(prev)
}

// === Restore ===
func RestoreKategori(
	prev *Kategori,
) common.ResultValue[*Kategori] {

	if prev == nil {
		return common.FailureValue[*Kategori](EmptyData())
	}

	prev.DeletedAt = nil

	return common.SuccessValue(prev)
}

// === Copy ===
func CopyKategori(
	prev *Kategori,
	copyCount int,
	createdby string,
	createdbyref string,
) common.ResultValue[*Kategori] {

	if prev == nil {
		return common.FailureValue[*Kategori](EmptyData())
	}

	var nama string

	if copyCount == 0 {
		nama = fmt.Sprintf("salin - %s", prev.NamaKategori)
	} else {
		nama = fmt.Sprintf("salin (%d) - %s", copyCount+1, prev.NamaKategori)
	}

	entity := &Kategori{
		UUID:         uuid.New(),
		NamaKategori: nama,
		SubKategori:  prev.SubKategori,
		CreatedBy:    helper.StrPtr(createdby),
		CreatedByRef: helper.StrPtr(createdbyref),
	}

	return common.SuccessValue(entity)
}

func isOverlap(start1, end1 time.Time) bool {
	return !end1.After(start1)
}
