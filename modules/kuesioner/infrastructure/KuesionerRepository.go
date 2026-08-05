package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	domainkuesioner "UnpakSiamida/modules/kuesioner/domain"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type KuesionerRepository struct {
	db *gorm.DB
}

func NewKuesionerRepository(db *gorm.DB) domainkuesioner.IKuesionerRepository {
	return &KuesionerRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *KuesionerRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainkuesioner.Kuesioner, error) {
	var Kuesioner domainkuesioner.Kuesioner

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&Kuesioner).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &Kuesioner, nil
}

// ------------------------
// GET DEFAULT BY UUID
// ------------------------
func (r *KuesionerRepository) GetDefaultByUuid(
	ctx context.Context,
	id uuid.UUID,
) (*domainkuesioner.KuesionerDefault, error) {

	var rowData domainkuesioner.KuesionerDefault

	err := r.db.WithContext(ctx).
		Table("kuesionerv2 a").
		Select(`
			a.id AS Id,
			a.uuid AS UUID,
			a.nidn AS NIDN,
			a.nama_dosen AS NamaDosen,
			a.nip AS NIP,
			a.nama_tendik AS NamaTendik,
			a.npm AS NPM,
			a.nama_mahasiswa AS NamaMahasiswa,
			a.kode_fakultas AS KodeFakultas,
			a.fakultas AS Fakultas,
			a.kode_prodi AS KodeProdi,
			a.prodi AS Prodi,
			a.unit AS Unit,
			a.id_bank_soal AS IdBankSoal,
			b.judul AS Judul,
			b.semester AS Semester,
			a.tanggal AS Tanggal
	`).
		Joins("JOIN bank_soalv2 b ON a.id_bank_soal = b.id").
		Where("a.uuid = ?", id).
		Take(&rowData).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &rowData, nil
}

func (r *KuesionerRepository) buildWhere(
	ctx context.Context,
	db *gorm.DB,
	JudulBankSoal *string,
	// Semester *string,
	Is4Year bool,
	partition_key string,
) (*gorm.DB, error) {

	// =====================================================
	// MODE 4 TAHUN
	// =====================================================
	if Is4Year {
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			return nil, err
		}

		now := time.Now().In(loc)
		start := time.Date(now.Year()-4, 1, 1, 0, 0, 0, 0, loc)
		end := time.Date(now.Year(), 12, 31, 0, 0, 0, 0, loc)
		startStr := start.Format("2006-01-02")
		endStr := end.Format("2006-01-02")

		db = db.Where("k.tanggal BETWEEN ? AND ?", startStr, endStr).Where("partition_key = ?", partition_key)
		// if partition_key == "UNIT" {
		// 	db = db.Where("k.tanggal BETWEEN ? AND ?", startStr, endStr).
		// 		Where("k.unit = ?", partition_key)
		// } else {
		// 	db = db.Where("k.tanggal BETWEEN ? AND ?", startStr, endStr).
		// 		Where("k.kode_fakultas = ?", partition_key)
		// }

		return db, nil
	}

	// =====================================================
	// MODE NORMAL (WAJIB JUDUL)
	// =====================================================
	rawStr := strings.TrimSpace(helper.NullableString(JudulBankSoal))
	if rawStr == "" {
		return nil, errors.New("judul bank soal wajib diisi")
	}

	// Parsing string "Judul (Semester)" contoh "Monev ... (202501)"
	lastOpen := strings.LastIndex(rawStr, " (")
	lastClose := strings.LastIndex(rawStr, ")")

	db = db.Where("partition_key = ?", partition_key)

	if lastOpen != -1 && lastClose > lastOpen {
		judul := strings.TrimSpace(rawStr[:lastOpen])
		semester := strings.TrimSpace(rawStr[lastOpen+2 : lastClose])
		db = db.Where("k.judul = ? AND k.semester = ?", judul, semester)
	} else {
		db = db.Where("k.judul = ?", rawStr)
	}

	return db, nil
}

func (r *KuesionerRepository) GetAllKuesionerResult(
	ctx context.Context,
	JudulBankSoal *string,
	// Semester *string,
	Is4Year bool,
	PartitionKey string,
) ([]domainkuesioner.KuesionerResult, error) {

	result := make([]domainkuesioner.KuesionerResult, 0)

	// =========================
	// BASE QUERY (JANGAN HILANG)
	// =========================
	db := r.db.WithContext(ctx).
		Table("kuesioner_materialized k")
	// 	Table("kuesionerv2 k").
	// 	Select(`
	// 	k.id,
	// 	k.uuid,
	// 	k.tanggal,

	// 	k.nidn,
	// 	k.nama_dosen,
	// 	k.nip,
	// 	k.nama_tendik,

	// 	k.npm,
	// 	k.nama_mahasiswa,

	// 	k.kode_fakultas,
	// 	k.fakultas,

	// 	k.kode_prodi,
	// 	k.prodi,
	// 	k.unit,

	// 	b.judul,
	// 	b.semester,

	// 	tp.pertanyaan,
	// 	tj.jawaban,
	// 	kj.freeText,

	// 	tp.jenis_pilihan,
	// 	ka.nama_kategori,
	// 	ka.full_text,

	// 	CASE
	// 		WHEN k.unit IS NOT NULL AND k.unit != '' THEN 'unit'
	// 		WHEN k.kode_fakultas IS NOT NULL AND k.kode_fakultas != '' THEN k.kode_fakultas
	// 		ELSE 'UNKNOWN'
	// 	END AS partition_key
	// `).
	// 	Joins("JOIN bank_soalv2 b ON b.id = k.id_bank_soal").
	// 	Joins("JOIN kuesioner_jawabanv2 kj ON kj.id_kuesioner = k.id").
	// 	Joins("LEFT JOIN template_pertanyaanv2 tp ON tp.id = kj.id_template_pertanyaan").
	// 	Joins("LEFT JOIN kategoriv2 ka ON ka.id = tp.id_kategori").
	// 	Joins("LEFT JOIN template_pilihanv2 tj ON tj.id = kj.id_template_jawaban")

	// Joins("LEFT JOIN m_dosen md ON md.NIDN = k.nidn").
	// Joins("LEFT JOIN users us1 ON us1.id = k.npm").
	// Joins("LEFT JOIN users us2 ON us2.id = k.nip").
	// Order("b.semester DESC, tp.pertanyaan ASC, tj.jawaban ASC")

	// =========================
	// APPLY WHERE (FIXED)
	// =========================
	db, err := r.buildWhere(ctx, db, JudulBankSoal, Is4Year, PartitionKey)
	if err != nil {
		return nil, err
	}

	// =========================
	// EXECUTE
	// =========================
	err = db.Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	"judul":          "b.judul",
	"semester":       "b.semester",
	"nidn":           "a.nidn",
	"nama_dosen":     "a.nama_dosen",
	"nip":            "a.nip",
	"nama_tendik":    "a.nama_tendik",
	"npm":            "a.npm",
	"nama_mahasiswa": "a.nama_mahasiswa",
	"fakultas":       "a.fakultas",
	"prodi":          "a.prodi",
	"unit":           "a.unit",
}

// ------------------------
// GET ALL
// ------------------------
func (r *KuesionerRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomain.SearchFilter,
	page, limit *int,
	deleted bool,
) ([]domainkuesioner.KuesionerDefault, int64, error) {

	var rows = make([]domainkuesioner.KuesionerDefault, 0)
	var total int64

	db := r.db.Debug().WithContext(ctx).
		Table("kuesionerv2 a").
		Select(`
			a.id AS Id,
			a.uuid AS UUID,
			a.nidn AS NIDN,
			a.nama_dosen AS NamaDosen,
			a.nip AS NIP,
			a.nama_tendik AS NamaTendik,
			a.npm AS NPM,
			a.nama_mahasiswa AS NamaMahasiswa,
			a.kode_fakultas AS KodeFakultas,
			a.fakultas AS Fakultas,
			a.kode_prodi AS KodeProdi,
			a.prodi AS Prodi,
			a.unit AS Unit,
			a.id_bank_soal AS IdBankSoal,
			b.judul AS Judul,
			b.semester AS Semester,
			a.tanggal AS Tanggal
	`).
		Joins("JOIN bank_soalv2 b ON a.id_bank_soal = b.id")

	if deleted {
		db = db.Where(clause.Expr{
			SQL: "b.deleted_at IS NOT NULL",
		})
	} else {
		db = db.Where(clause.Expr{
			SQL: "b.deleted_at IS NULL",
		})
	}

	// -----------------------------------
	// ADVANCED FILTERS
	// -----------------------------------
	for _, f := range searchFilters {
		col, ok := allowedSearchColumns[strings.ToLower(f.Field)]
		if !ok {
			continue
		}

		val := ""
		if f.Value != nil {
			val = strings.TrimSpace(*f.Value)
		}
		if val == "" {
			continue
		}

		switch strings.ToLower(f.Operator) {
		case "eq":
			db = db.Where(clause.Eq{
				Column: col,
				Value:  val,
			})
		case "neq":
			db = db.Where(clause.Neq{
				Column: col,
				Value:  val,
			})
		case "like":
			db = db.Where(clause.Like{
				Column: col,
				Value:  "%" + helper.EscapeLike(val) + "%",
			})
		case "in":
			rawVals := strings.Split(val, ",")
			vals := make([]interface{}, 0, len(rawVals))

			for _, v := range rawVals {
				v = strings.TrimSpace(v)
				if v != "" {
					vals = append(vals, v)
				}
			}

			if len(vals) > 0 {
				db = db.Where(clause.IN{
					Column: col,
					Values: vals,
				})
			}
		}
	}

	// -----------------------------------
	// GLOBAL SEARCH
	// -----------------------------------
	if strings.TrimSpace(search) != "" {
		like := "%" + search + "%"
		var conditions []clause.Expression

		for _, col := range allowedSearchColumns {
			conditions = append(conditions, clause.Like{
				Column: col,
				Value:  like,
			})
		}

		if len(conditions) > 0 {
			db = db.Where(clause.Or(conditions...))
		}
	}

	// -----------------------------------
	// COUNT (AMAN)
	// -----------------------------------
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// -----------------------------------
	// ORDER + PAGINATION
	// -----------------------------------
	db = db.Order("a.id DESC")

	if page != nil && limit != nil && *limit > 0 {
		offset := (*page - 1) * (*limit)
		db = db.Offset(offset).Limit(*limit)
	}

	// -----------------------------------
	// EXECUTE
	// -----------------------------------
	if err := db.Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *KuesionerRepository) GetAllFormFromActiveBankSoal(
	ctx context.Context,
	nidn string,
	nip string,
	npm string,
	banksoal []uint,
) ([]domainkuesioner.KuesionerDefault, error) {

	var rows = make([]domainkuesioner.KuesionerDefault, 0)

	db := r.db.Debug().WithContext(ctx).
		Table("kuesionerv2 a").
		Select(`
			a.id AS Id,
			a.uuid AS UUID,
			a.nidn AS NIDN,
			a.nama_dosen AS NamaDosen,
			a.nip AS NIP,
			a.nama_tendik AS NamaTendik,
			a.npm AS NPM,
			a.nama_mahasiswa AS NamaMahasiswa,
			a.kode_fakultas AS KodeFakultas,
			a.fakultas AS Fakultas,
			a.kode_prodi AS KodeProdi,
			a.prodi AS Prodi,
			a.unit AS Unit,
			a.id_bank_soal AS IdBankSoal,
			b.judul AS Judul,
			b.semester AS Semester,
			a.tanggal AS Tanggal
	`).
		Joins("JOIN bank_soalv2 b ON a.id_bank_soal = b.id").
		Where("a.nidn = ?", nidn).
		Where("a.nip = ?", nip).
		Where("a.npm = ?", npm).
		Where("a.id_bank_soal in (?)", banksoal)

	// -----------------------------------
	// ORDER
	// -----------------------------------
	db = db.Order("a.id DESC")

	// -----------------------------------
	// EXECUTE
	// -----------------------------------
	if err := db.Find(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

// ------------------------
// CREATE
// ------------------------
func (r *KuesionerRepository) Create(ctx context.Context, kuesioner *domainkuesioner.Kuesioner) error {
	tx := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id_bank_soal"},
				{Name: "nidn"},
				{Name: "nip"},
				{Name: "npm"},
			},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"id":             gorm.Expr("LAST_INSERT_ID(id)"),
				"nama_dosen":     kuesioner.NamaDosen,
				"nama_tendik":    kuesioner.NamaTendik,
				"nama_mahasiswa": kuesioner.NamaMahasiswa,
				"kode_fakultas":  kuesioner.KodeFakultas,
				"fakultas":       kuesioner.Fakultas,
				"kode_prodi":     kuesioner.KodeProdi,
				"prodi":          kuesioner.Prodi,
				"unit":           kuesioner.Unit,
				"tanggal":        kuesioner.Tanggal,
				"createdBy":      kuesioner.CreatedBy,
				"createdByRef":   kuesioner.CreatedByRef,
				"updated_at":     gorm.Expr("NOW()"),
			}),
		}).
		Create(kuesioner)

	if tx.Error != nil {
		return tx.Error
	}

	var result struct {
		ID   uint
		UUID uuid.UUID
	}

	if err := r.db.Raw(`
		SELECT id, uuid 
		FROM kuesionerv2 
		WHERE id = LAST_INSERT_ID()
	`).Scan(&result).Error; err != nil {
		return err
	}

	kuesioner.ID = result.ID
	kuesioner.UUID = result.UUID

	return nil
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *KuesionerRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domainkuesioner.Kuesioner{}).Error
}

func (r *KuesionerRepository) SetupUuid(ctx context.Context) error {
	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domainkuesioner.Kuesioner{}).
		Where("uuid IS NULL OR uuid = ''").
		Pluck("id", &ids).Error; err != nil {
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	for i := 0; i < len(ids); i += chunkSize {
		end := i + chunkSize
		if end > len(ids) {
			end = len(ids)
		}

		chunk := ids[i:end]

		caseSQL := "CASE id "
		args := make([]any, 0, len(chunk)*2+1)

		for _, id := range chunk {
			u := uuid.NewString()
			caseSQL += "WHEN ? THEN ? "
			args = append(args, id, u)
		}

		caseSQL += "END"
		args = append(args, chunk)

		query := fmt.Sprintf(
			"UPDATE kuesionerv2 SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}

func applySummaryFilter(db *gorm.DB, rawJudul string) *gorm.DB {
	rawStr := strings.TrimSpace(rawJudul)
	var judul, semester string
	lastOpen := strings.LastIndex(rawStr, " (")
	lastClose := strings.LastIndex(rawStr, ")")
	if lastOpen != -1 && lastClose > lastOpen {
		judul = strings.TrimSpace(rawStr[:lastOpen])
		semester = strings.TrimSpace(rawStr[lastOpen+2 : lastClose])
	} else {
		judul = rawStr
	}

	cleanJudul := strings.TrimSpace(
		strings.ReplaceAll(
			strings.ReplaceAll(judul, " - Semester Ganjil", ""),
			" - Semester Genap", "",
		),
	)

	if semester != "" {
		return db.Where("(judul LIKE ? OR judul LIKE ?) AND semester = ?", "%"+judul+"%", "%"+cleanJudul+"%", semester)
	}
	return db.Where("judul LIKE ?", "%"+cleanJudul+"%")
}

func (r *KuesionerRepository) GetReportSummaryOverview(ctx context.Context, rawJudul, kodeFakultas, kodeProdi string) (*domainkuesioner.ReportSummaryOverview, error) {
	var overview domainkuesioner.ReportSummaryOverview
	qOverview := applySummaryFilter(r.db.WithContext(ctx).Table("report_summary_overview"), rawJudul).
		Where("kode_fakultas = ? AND kode_prodi = ?", strings.TrimSpace(kodeFakultas), strings.TrimSpace(kodeProdi))
	err := qOverview.First(&overview).Error
	return &overview, err
}

func (r *KuesionerRepository) GetReportDistribusiFakultas(ctx context.Context, rawJudul, kodeFakultas, kodeProdi string) ([]domainkuesioner.ReportDistribusiFakultas, error) {
	var fakultas []domainkuesioner.ReportDistribusiFakultas
	qFakultas := applySummaryFilter(r.db.WithContext(ctx).Table("report_distribusi_fakultas"), rawJudul).
		Where("kode_fakultas = ? AND kode_prodi = ?", strings.TrimSpace(kodeFakultas), strings.TrimSpace(kodeProdi))
	err := qFakultas.Order("total_responden DESC").Find(&fakultas).Error
	return fakultas, err
}

func (r *KuesionerRepository) GetReportTopQuestions(ctx context.Context, rawJudul, kodeFakultas, kodeProdi string) ([]domainkuesioner.ReportTopQuestion, error) {
	var topQuestions []domainkuesioner.ReportTopQuestion
	qTop := applySummaryFilter(r.db.WithContext(ctx).Table("report_top_questions"), rawJudul).
		Where("kode_fakultas = ? AND kode_prodi = ?", strings.TrimSpace(kodeFakultas), strings.TrimSpace(kodeProdi))
	err := qTop.Order("peringkat ASC").Limit(10).Find(&topQuestions).Error
	return topQuestions, err
}

func (r *KuesionerRepository) GetReportKategoriSummary(ctx context.Context, rawJudul, kodeFakultas, kodeProdi string) ([]domainkuesioner.ReportKategoriSummary, error) {
	var kategoriSummary []domainkuesioner.ReportKategoriSummary
	qKat := applySummaryFilter(r.db.WithContext(ctx).Table("report_kategori_summary"), rawJudul).
		Where("kode_fakultas = ? AND kode_prodi = ?", strings.TrimSpace(kodeFakultas), strings.TrimSpace(kodeProdi))
	err := qKat.Order("nama_kategori ASC").Find(&kategoriSummary).Error
	return kategoriSummary, err
}

func (r *KuesionerRepository) GetReportYear(ctx context.Context) ([]domainkuesioner.ReportYear, error) {
	var list []domainkuesioner.ReportYear
	err := r.db.WithContext(ctx).Table("report_year").Order("tahun ASC").Find(&list).Error
	return list, err
}

func (r *KuesionerRepository) GetDashboardStats(ctx context.Context) (*domainkuesioner.DashboardStats, error) {
	var stats domainkuesioner.DashboardStats
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			(SELECT COALESCE(SUM(total_responden), 0) FROM report_summary_overview) AS total_responden,
			(SELECT COUNT(*) FROM bank_soalv2 WHERE deleted_at IS NULL) AS active_surveys,
			(SELECT COUNT(DISTINCT kode_prodi) FROM m_program_studi) AS total_prodi,
			(SELECT COUNT(DISTINCT kode_fakultas) FROM m_fakultas) AS total_fakultas,
			(SELECT COALESCE(ROUND(AVG(rata_rata_rating), 2), 0.00) FROM report_summary_overview WHERE rata_rata_rating > 0) AS rata_rata_rating
	`).Scan(&stats).Error
	return &stats, err
}

func (r *KuesionerRepository) GetReportSummary(ctx context.Context, rawJudul, kodeFakultas, kodeProdi string) (*domainkuesioner.ReportSummaryResponse, error) {
	overview, _ := r.GetReportSummaryOverview(ctx, rawJudul, kodeFakultas, kodeProdi)
	fakultas, _ := r.GetReportDistribusiFakultas(ctx, rawJudul, kodeFakultas, kodeProdi)
	topQuestions, _ := r.GetReportTopQuestions(ctx, rawJudul, kodeFakultas, kodeProdi)
	kategoriSummary, _ := r.GetReportKategoriSummary(ctx, rawJudul, kodeFakultas, kodeProdi)
	reportYear, _ := r.GetReportYear(ctx)

	return &domainkuesioner.ReportSummaryResponse{
		Overview:           overview,
		DistribusiFakultas: fakultas,
		TopQuestions:       topQuestions,
		KategoriSummary:    kategoriSummary,
		ReportYear:         reportYear,
	}, nil
}
