package application

type UpdateBankSoalCommand struct {
	Uuid         string
	Judul        string
	Content      string
	Deskripsi    string
	Semester     string
	TanggalMulai string
	TanggalAkhir string
	SID          string
	Resource     string
}
