package domain

import (
	"time"

	"github.com/google/uuid"
)

type KuesionerDefault struct {
	Id            uint
	UUID          uuid.UUID
	NIDN          *string
	NamaDosen     *string
	NIP           *string
	NamaTendik    *string
	NPM           *string
	NamaMahasiswa *string
	KodeFakultas  *string
	Fakultas      *string
	KodeProdi     *string
	Prodi         *string
	Unit          *string
	IdBankSoal    string
	UUIDBankSoal  uuid.UUID
	Judul         string
	Semester      *string
	Tanggal       time.Time
}
