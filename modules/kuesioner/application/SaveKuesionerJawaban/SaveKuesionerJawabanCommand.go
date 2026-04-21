package application

type SaveKuesionerJawabanCommand struct {
	UuidKuesioner  string
	UuidPertanyaan string
	Jawaban        string
	SID            string
	Resource       string
	CodeCtx        string
}
