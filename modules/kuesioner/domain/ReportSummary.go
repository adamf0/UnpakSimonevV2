package domain

type ReportSummaryOverview struct {
	ID             uint    `gorm:"primaryKey;column:id" json:"id"`
	Judul          string  `gorm:"column:judul" json:"judul"`
	Semester       string  `gorm:"column:semester" json:"semester"`
	TotalResponden int     `gorm:"column:total_responden" json:"total_responden"`
	TotalJawaban   int     `gorm:"column:total_jawaban" json:"total_jawaban"`
	RataRataRating float64 `gorm:"column:rata_rata_rating" json:"rata_rata_rating"`
}

func (ReportSummaryOverview) TableName() string {
	return "report_summary_overview"
}

type ReportDistribusiFakultas struct {
	ID                uint    `gorm:"primaryKey;column:id" json:"id"`
	Judul             string  `gorm:"column:judul" json:"judul"`
	Semester          string  `gorm:"column:semester" json:"semester"`
	KodeFakultas      string  `gorm:"column:kode_fakultas" json:"kode_fakultas"`
	NamaFakultas      string  `gorm:"column:nama_fakultas" json:"nama_fakultas"`
	TotalResponden    int     `gorm:"column:total_responden" json:"total_responden"`
	Persentase        float64 `gorm:"column:persentase" json:"persentase"`
	ProdiDistribution string  `gorm:"column:prodi_distribution;type:json" json:"prodi_distribution"`
}

func (ReportDistribusiFakultas) TableName() string {
	return "report_distribusi_fakultas"
}

type ReportTopQuestion struct {
	ID             uint    `gorm:"primaryKey;column:id" json:"id"`
	Judul          string  `gorm:"column:judul" json:"judul"`
	Semester       string  `gorm:"column:semester" json:"semester"`
	Pertanyaan     string  `gorm:"column:pertanyaan" json:"pertanyaan"`
	NamaKategori   string  `gorm:"column:nama_kategori" json:"nama_kategori"`
	TotalResponden int     `gorm:"column:total_responden" json:"total_responden"`
	RataRataSkor   float64 `gorm:"column:rata_rata_skor" json:"rata_rata_skor"`
	Peringkat      int     `gorm:"column:peringkat" json:"peringkat"`
}

func (ReportTopQuestion) TableName() string {
	return "report_top_questions"
}

type ReportKategoriSummary struct {
	ID                uint    `gorm:"primaryKey;column:id" json:"id"`
	Judul             string  `gorm:"column:judul" json:"judul"`
	Semester          string  `gorm:"column:semester" json:"semester"`
	NamaKategori      string  `gorm:"column:nama_kategori" json:"nama_kategori"`
	FullText          string  `gorm:"column:full_text" json:"full_text"`
	TotalPertanyaan   int     `gorm:"column:total_pertanyaan" json:"total_pertanyaan"`
	TotalResponden    int     `gorm:"column:total_responden" json:"total_responden"`
	RataRataSkor      float64 `gorm:"column:rata_rata_skor" json:"rata_rata_skor"`
	ChartDistribution string  `gorm:"column:chart_distribution;type:json" json:"chart_distribution"`
	QuestionsJson     string  `gorm:"column:questions_json;type:json" json:"questions_json"`
}

func (ReportKategoriSummary) TableName() string {
	return "report_kategori_summary"
}

type ReportYear struct {
	ID             uint   `gorm:"primaryKey;column:id" json:"id"`
	Tahun          string `gorm:"column:tahun" json:"tahun"`
	TotalKuesioner int    `gorm:"column:total_kuesioner" json:"total_kuesioner"`
	TotalMahasiswa int    `gorm:"column:total_mahasiswa" json:"total_mahasiswa"`
	TotalDosen     int    `gorm:"column:total_dosen" json:"total_dosen"`
	TotalTendik    int    `gorm:"column:total_tendik" json:"total_tendik"`
}

func (ReportYear) TableName() string {
	return "report_year"
}

type ReportSummaryResponse struct {
	Overview           *ReportSummaryOverview     `json:"overview"`
	DistribusiFakultas []ReportDistribusiFakultas `json:"distribusi_fakultas"`
	TopQuestions       []ReportTopQuestion        `json:"top_questions"`
	KategoriSummary    []ReportKategoriSummary    `json:"kategori_summary"`
	ReportYear         []ReportYear               `json:"report_year"`
}
