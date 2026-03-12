package domain

import (
	"time"

	common "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"

	"github.com/google/uuid"
)

type Kuesioner struct {
	common.Entity

	ID            uint      `gorm:"primaryKey;autoIncrement"`
	UUID          uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	NIDN          string    `gorm:"column:nidn"`
	NamaDosen     *string   `gorm:"column:nama_dosen"`
	NIP           string    `gorm:"column:nip"`
	NamaTendik    *string   `gorm:"column:nama_tendik"`
	NPM           string    `gorm:"column:npm"`
	NamaMahasiswa *string   `gorm:"column:nama_mahasiswa"`
	KodeFakultas  *string   `gorm:"column:kode_fakultas"`
	Fakultas      *string   `gorm:"column:fakultas"`
	KodeProdi     *string   `gorm:"column:kode_prodi"`
	Prodi         *string   `gorm:"column:prodi"`
	Unit          *string   `gorm:"column:unit"`
	IdBankSoal    string    `gorm:"column:id_bank_soal"`
	Tanggal       time.Time `gorm:"column:tanggal"`

	CreatedBy    *string `gorm:"column:createdBy"`
	CreatedByRef *string `gorm:"column:createdByRef"`
}

func (Kuesioner) TableName() string {
	return "kuesionerv2"
}

// === CREATE ===
func NewKuesioner(
	nidn *string,
	namaDosen *string,
	nip *string,
	namaTendik *string,
	npm *string,
	namaMahasiswa *string,
	kodeFakultas *string,
	fakultas *string,
	kodeProdi *string,
	prodi *string,
	unit *string,
	idBankSoal string,
	tanggal string,
	createdby string, //lpm, fakultas, prodi
	createdbyref string,
) common.ResultValue[*Kuesioner] {
	format := "2006-01-02 15:04:05"

	tanggalParse, err := time.Parse(format, tanggal)
	if err != nil {
		return common.FailureValue[*Kuesioner](InvalidDate("tanggal mulai"))
	}

	kuesioner := &Kuesioner{
		UUID:          uuid.New(),
		NIDN:          helper.StringValue(nidn),
		NamaDosen:     namaDosen,
		NIP:           helper.StringValue(nip),
		NamaTendik:    namaTendik,
		NPM:           helper.StringValue(npm),
		NamaMahasiswa: namaMahasiswa,
		KodeFakultas:  kodeFakultas,
		Fakultas:      fakultas,
		KodeProdi:     kodeProdi,
		Prodi:         prodi,
		Unit:          unit,
		IdBankSoal:    idBankSoal,
		Tanggal:       tanggalParse,
		CreatedBy:     helper.StrPtr(createdby),
		CreatedByRef:  helper.StrPtr(createdbyref),
	}

	return common.SuccessValue(kuesioner)
}
