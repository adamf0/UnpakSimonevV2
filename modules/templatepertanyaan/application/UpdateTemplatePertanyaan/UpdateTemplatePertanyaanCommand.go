package application

type UpdateTemplatePertanyaanCommand struct {
	Uuid         string
	UuidBankSoal string
	Pertanyaan   string
	JenisPilihan string
	Bobot        string
	UuidKategori string
	Required     int
	SID          string
	Resource     string
}
