package application

type UpdateKategoriCommand struct {
	Uuid         string
	NamaKategori string
	SubKategori  *string
	SID          string
	Resource     string
}
