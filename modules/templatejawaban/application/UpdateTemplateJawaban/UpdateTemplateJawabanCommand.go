package application

type UpdateTemplateJawabanCommand struct {
	Uuid                   string
	UuidTemplatePertanyaan string
	Jawaban                string
	Nilai                  string
	IsFreeText             string
	SID                    string
	Resource               string
}
