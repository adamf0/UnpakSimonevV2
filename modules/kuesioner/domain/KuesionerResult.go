package domain

type KuesionerResult struct {
	ID            uint `json:"-"`
	UUID          string
	NIDN          string
	NamaDosen     string
	NIP           string
	NamaTendik    string
	NPM           string
	NamaMahasiswa string

	KodeFakultas string
	Fakultas     string
	KodeProdi    string
	Prodi        string
	Unit         string

	Judul        string
	Semester     int
	Pertanyaan   string
	Jawaban      string
	FreeText     string
	JenisPilihan string
	Kategori     string
	FullPath     string
}
