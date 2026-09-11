package SaveBulkKuesionerJawaban

type BulkJawabanItem struct {
	UuidPertanyaan string `json:"pertanyaan"`
	Jawaban        string `json:"jawaban"`
}

type SaveBulkKuesionerJawabanCommand struct {
	UuidKuesioner string
	SID           string
	Resource      string
	CodeCtx       string
	Items         []BulkJawabanItem
}
