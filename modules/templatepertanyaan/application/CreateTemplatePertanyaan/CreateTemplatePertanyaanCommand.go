package application

type CreateTemplatePertanyaanCommand struct {
	UuidBankSoal string
	Pertanyaan   string
	JenisPilihan string
	Bobot        string
	UuidKategori string
	Required     int
	SID          string
	Resource     string
}
