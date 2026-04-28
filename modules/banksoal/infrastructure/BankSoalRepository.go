package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	domainbanksoal "UnpakSiamida/modules/banksoal/domain"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BankSoalRepository struct {
	db *gorm.DB
}

func NewBankSoalRepository(db *gorm.DB) domainbanksoal.IBankSoalRepository {
	return &BankSoalRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *BankSoalRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainbanksoal.BankSoal, error) {
	var BankSoal domainbanksoal.BankSoal

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&BankSoal).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &BankSoal, nil
}

// ------------------------
// GET DEFAULT BY UUID
// ------------------------
// func (r *BankSoalRepository) GetDefaultByUuid( //[note] ini lebih optimal dibandingkan getall ketika di analyze
// 	ctx context.Context,
// 	id uuid.UUID,
// ) (*domainbanksoal.BankSoalDefault, error) {

// 	var row domainbanksoal.BankSoalDefault

// 	db := r.db.WithContext(ctx).
// 		Table("bank_soalv2 b").
// 		Debug().
// 		Select(`
// 		b.id AS ID,
// 		b.uuid AS UUID,
// 		b.judul AS Judul,
// 		b.content AS Content,
// 		b.deskripsi AS Deskripsi,
// 		b.semester AS Semester,
// 		b.tanggal_mulai AS TanggalMulai,
// 		b.tanggal_akhir AS TanggalAkhir,
// 		b.createdBy AS CreatedBy,
// 		b.createdByRef AS CreatedByRef,
// 		b.deleted_at AS DeletedAt,
// 		b.status AS Status,

// 		COALESCE(u.name, d.nama_dosen) AS Nama,
// 		COALESCE(u.level, 'dosen') AS Role,
// 		COALESCE(fu.kode_fakultas, fd.kode_fakultas) AS KodeFakultas,
// 		COALESCE(fu.nama_fakultas, fd.nama_fakultas) AS NamaFakultas,
// 		COALESCE(pu.kode_prodi, pd.kode_prodi) AS KodeProdi,
// 		COALESCE(pu.nama_prodi, pd.nama_prodi) AS NamaProdi,

// 		COALESCE(pc.total_pertanyaan, 0) AS TotalPertanyaan,
// 		0 AS TotalInput,
// 		k.uuid AS UUIDKuesioner
// 	`).
// 		Joins(`LEFT JOIN users u ON u.id = b.createdByRef AND b.createdBy = 'local'`).
// 		Joins(`LEFT JOIN m_fakultas fu ON fu.kode_fakultas = u.fakultas`).
// 		Joins(`LEFT JOIN m_program_studi pu ON pu.kode_prodi = u.prodi`).
// 		Joins(`LEFT JOIN m_dosen d ON d.NIDN = b.createdByRef AND b.createdBy = 'simak'`).
// 		Joins(`LEFT JOIN m_fakultas fd ON fd.kode_fakultas = d.kode_fak`).
// 		Joins(`LEFT JOIN m_program_studi pd ON pd.kode_prodi = d.kode_prodi`).
// 		Joins(`
// 		LEFT JOIN (
// 			SELECT id_bank_soal, COUNT(id) AS total_pertanyaan
// 			FROM template_pertanyaanv2
// 			GROUP BY id_bank_soal
// 		) pc ON pc.id_bank_soal = b.id
// 	`).
// 		Joins(`LEFT JOIN kuesionerv2 k ON k.id_bank_soal = b.id`).
// 		Where("b.uuid = ?", id).
// 		Order("b.id DESC")

// 	if err := db.First(&row).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, gorm.ErrRecordNotFound
// 		}
// 		return nil, err
// 	}

//		return &row, nil
//	}
func (r *BankSoalRepository) GetDefaultByUuid(
	ctx context.Context,
	id uuid.UUID,
) (*domainbanksoal.BankSoalDefault, error) {

	var row domainbanksoal.BankSoalDefault

	// =========================
	// Subquery dosen
	// =========================
	dosenSub := r.db.
		Table("m_dosen d").
		Select(`
			CAST(d.NIDN AS CHAR) AS nidn,
			d.nama_dosen,
			f.kode_fakultas,
			f.nama_fakultas,
			p.kode_prodi,
			p.kode_jenjang,
			CONCAT(
				p.nama_prodi,
				CASE p.kode_jenjang
					WHEN 'J' THEN ' (Profesi)'
					WHEN 'E' THEN ' (D3)'
					WHEN 'D' THEN ' (D4)'
					WHEN 'A' THEN ' (S3)'
					WHEN 'B' THEN ' (S2)'
					WHEN 'C' THEN ' (S1)'
					ELSE ''
				END
			) AS nama_prodi,
			'dosen' AS role
		`).
		Joins("LEFT JOIN m_fakultas f ON d.kode_fak = f.kode_fakultas").
		Joins("LEFT JOIN m_program_studi p ON d.kode_prodi = p.kode_prodi")

	// =========================
	// Subquery account
	// =========================
	accountSub := r.db.
		Table("users u").
		Select(`
			CAST(u.id AS CHAR) AS id,
			u.name,
			f.kode_fakultas,
			f.nama_fakultas,
			p.kode_prodi,
			p.kode_jenjang,
			CONCAT(
				p.nama_prodi,
				CASE p.kode_jenjang
					WHEN 'J' THEN ' (Profesi)'
					WHEN 'E' THEN ' (D3)'
					WHEN 'D' THEN ' (D4)'
					WHEN 'A' THEN ' (S3)'
					WHEN 'B' THEN ' (S2)'
					WHEN 'C' THEN ' (S1)'
					ELSE ''
				END
			) AS nama_prodi,
			u.level AS role
		`).
		Joins("LEFT JOIN m_fakultas f ON u.fakultas = f.kode_fakultas").
		Joins("LEFT JOIN m_program_studi p ON u.prodi = p.kode_prodi")

	// =========================
	// Subquery pertanyaan
	// =========================
	pertanyaanSub := r.db.
		Table("template_pertanyaanv2 tp").
		Select(`
			tp.id_bank_soal,
			COUNT(tp.id) AS total_pertanyaan,
			GROUP_CONCAT(tp.uuid) AS uuids
		`).
		Where("tp.deleted_at IS NULL").
		Where("tp.status = 'active'").
		Group("tp.id_bank_soal")

	// =========================
	// Main query
	// =========================
	db := r.db.WithContext(ctx).
		Table("bank_soalv2 b").
		Debug().
		Select(`
			b.id AS ID,
			b.uuid AS UUID,
			b.judul AS Judul,
			b.content AS Content,
			b.deskripsi AS Deskripsi,
			b.semester AS Semester,
			b.tanggal_mulai AS TanggalMulai,
			b.tanggal_akhir AS TanggalAkhir,
			b.createdBy AS CreatedBy,
			b.createdByRef AS CreatedByRef,
			b.deleted_at AS DeletedAt,
			b.status AS Status,

			COALESCE(ul.name, dc.nama_dosen) AS Nama,
			COALESCE(ul.role, dc.role) AS Role,
			COALESCE(ul.kode_fakultas, dc.kode_fakultas) AS KodeFakultas,
			COALESCE(ul.nama_fakultas, dc.nama_fakultas) AS NamaFakultas,
			COALESCE(ul.kode_prodi, dc.kode_prodi) AS KodeProdi,
			COALESCE(ul.nama_prodi, dc.nama_prodi) AS NamaProdi,

			COALESCE(pc.total_pertanyaan, 0) AS TotalPertanyaan,
			0 AS TotalInput,
			COALESCE(pc.uuids, '') AS RawTargetUUIDs,
			k.uuid AS UUIDKuesioner
		`).
		Joins(`LEFT JOIN (?) ul 
			ON ul.id = CAST(b.createdByRef AS CHAR)
			AND LOWER(b.createdBy) = 'local'`, accountSub).
		Joins(`LEFT JOIN (?) dc 
			ON dc.nidn = CAST(b.createdByRef AS CHAR)
			AND LOWER(b.createdBy) = 'simak'`, dosenSub).
		Joins(`LEFT JOIN (?) pc ON b.id = pc.id_bank_soal`, pertanyaanSub).
		Joins(`LEFT JOIN kuesionerv2 k ON k.id_bank_soal = b.id`).
		Where("b.uuid = ?", id).
		Where("b.deleted_at IS NULL").
		Order("b.id DESC")

	if err := db.First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	// =========================
	// LOAD EXTENSION DATA
	// =========================
	var extRows []domainbanksoal.BankSoalExtDefault

	err := r.db.
		Table("bank_soal_extendv2 e").
		Debug().
		Select(`
			e.id AS ID,
			e.uuid AS UUID,
			e.id_bank_soal AS IdBankSoal,
			e.tanggal_mulai AS TanggalMulai,
			e.tanggal_akhir AS TanggalAkhir,
			e.createdBy AS CreatedBy,
			e.createdByRef AS CreatedByRef,

			COALESCE(ul.kode_fakultas, dc.kode_fakultas) AS KodeFakultas,
			COALESCE(ul.nama_fakultas, dc.nama_fakultas) AS NamaFakultas,
			COALESCE(ul.kode_prodi, dc.kode_prodi) AS KodeProdi,
			COALESCE(ul.nama_prodi, dc.nama_prodi) AS NamaProdi,
			COALESCE(ul.role, dc.role) AS Role
		`).
		Joins(`LEFT JOIN (?) ul 
			ON ul.id = CAST(e.createdByRef AS CHAR)
			AND LOWER(e.createdBy) = 'local'`, accountSub).
		Joins(`LEFT JOIN (?) dc 
			ON dc.nidn = CAST(e.createdByRef AS CHAR)
			AND LOWER(e.createdBy) = 'simak'`, dosenSub).
		Where("e.id_bank_soal = ?", row.Id).
		Order("e.id DESC").
		Scan(&extRows).Error

	if err != nil {
		return nil, err
	}

	row.ListExt = append(
		[]domainbanksoal.BankSoalExtDefault{},
		extRows...,
	)

	// =========================
	// Parse UUID Pertanyaan
	// =========================
	if row.RawTargetUUIDs != "" {
		uuids := strings.Split(row.RawTargetUUIDs, ",")
		for _, s := range uuids {
			if u, err := uuid.Parse(strings.TrimSpace(s)); err == nil {
				row.TargetPertanyaan = append(row.TargetPertanyaan, u)
			}
		}
	}

	return &row, nil
}

func (r *BankSoalRepository) GetDefaultByKuesioner(
	ctx context.Context,
	id uuid.UUID,
) (*domainbanksoal.BankSoalDefault, error) {

	var row domainbanksoal.BankSoalDefault

	// Subquery dosen
	dosenSub := r.db.
		Table("m_dosen d").
		Select(`
			CAST(d.NIDN AS CHAR) AS nidn,
			d.nama_dosen,
			f.kode_fakultas,
			f.nama_fakultas,
			p.kode_prodi,
			p.kode_jenjang,
			CONCAT(
				p.nama_prodi,
				CASE p.kode_jenjang
					WHEN 'J' THEN ' (Profesi)'
					WHEN 'E' THEN ' (D3)'
					WHEN 'D' THEN ' (D4)'
					WHEN 'A' THEN ' (S3)'
					WHEN 'B' THEN ' (S2)'
					WHEN 'C' THEN ' (S1)'
					ELSE ''
				END
			) AS nama_prodi,
			'dosen' as role
		`).
		Joins("LEFT JOIN m_fakultas f ON d.kode_fak = f.kode_fakultas").
		Joins("LEFT JOIN m_program_studi p ON d.kode_prodi = p.kode_prodi")

	// Subquery account
	accountSub := r.db.
		Table("users u").
		Select(`
			CAST(u.id AS CHAR) AS id,
			u.name,
			f.kode_fakultas,
			f.nama_fakultas,
			p.kode_prodi,
			p.kode_jenjang,
			CONCAT(
				p.nama_prodi,
				CASE p.kode_jenjang
					WHEN 'J' THEN ' (Profesi)'
					WHEN 'E' THEN ' (D3)'
					WHEN 'D' THEN ' (D4)'
					WHEN 'A' THEN ' (S3)'
					WHEN 'B' THEN ' (S2)'
					WHEN 'C' THEN ' (S1)'
					ELSE ''
				END
			) AS nama_prodi,
			u.level as role
		`).
		Joins("LEFT JOIN m_fakultas f ON u.fakultas = f.kode_fakultas").
		Joins("LEFT JOIN m_program_studi p ON u.prodi = p.kode_prodi")

	// =========================
	// SUBQUERY PERTANYAAN
	// =========================
	pertanyaanSub := r.db.
		Table("template_pertanyaanv2 tp").
		Select(`
			tp.id_bank_soal,
			COUNT(tp.id) as total_pertanyaan,
			GROUP_CONCAT(tp.uuid) as uuids
		`).
		Where("tp.deleted_at IS NULL").
		Where("tp.status = ?", "active").
		Group("tp.id_bank_soal")

	// =========================
	// MAIN QUERY
	// =========================
	db := r.db.WithContext(ctx).
		Table("bank_soalv2 b").
		Debug().
		Select(`
			b.id as ID,
			b.uuid as UUID,
			b.judul as Judul,
			b.content as Content,
			b.deskripsi as Deskripsi,
			b.semester as Semester,
			b.tanggal_mulai as TanggalMulai,
			b.tanggal_akhir as TanggalAkhir,
			b.createdBy as CreatedBy,
			b.createdByRef as CreatedByRef,
			b.deleted_at as DeletedAt,
			b.status as Status,

			k.uuid as UUIDKuesioner,

			COALESCE(ul.name, dc.nama_dosen) AS Nama,
			COALESCE(ul.role, dc.role) AS Role,
			COALESCE(ul.kode_fakultas, dc.kode_fakultas) AS KodeFakultas,
			COALESCE(ul.nama_fakultas, dc.nama_fakultas) AS NamaFakultas,
			COALESCE(ul.kode_prodi, dc.kode_prodi) AS KodeProdi,
			COALESCE(ul.nama_prodi, dc.nama_prodi) AS NamaProdi,

			COALESCE(pc.total_pertanyaan, 0) AS TotalPertanyaan,
			0 AS TotalInput,
			COALESCE(pc.uuids, '') AS TargetPertanyaan
		`).
		Joins(`
			INNER JOIN kuesionerv2 k
				ON k.id_bank_soal = b.id
		`).
		Joins(`
			LEFT JOIN (?) ul
				ON ul.id = CAST(b.createdByRef AS CHAR)
				AND LOWER(b.createdBy) = 'local'
		`, accountSub).
		Joins(`
			LEFT JOIN (?) dc
				ON dc.nidn = CAST(b.createdByRef AS CHAR)
				AND LOWER(b.createdBy) = 'simak'
		`, dosenSub).
		Joins(`
			LEFT JOIN (?) pc
				ON b.id = pc.id_bank_soal
		`, pertanyaanSub).
		Where("k.uuid = ?", id).
		Where("b.deleted_at IS NULL").
		Order("b.id DESC").
		Limit(1)

	// =========================
	// EXECUTE
	// =========================
	if err := db.First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	// =========================
	// LOAD EXTENSION DATA
	// =========================
	extRows := make([]domainbanksoal.BankSoalExtDefault, 0)

	err := r.db.WithContext(ctx).
		Table("bank_soal_extendv2 e").
		Select(`
			e.id as ID,
			e.uuid as UUID,
			e.id_bank_soal as IdBankSoal,
			e.tanggal_mulai as TanggalMulai,
			e.tanggal_akhir as TanggalAkhir,
			e.createdBy as CreatedBy,
			e.createdByRef as CreatedByRef,

			COALESCE(ul.kode_fakultas, dc.kode_fakultas) AS KodeFakultas,
			COALESCE(ul.nama_fakultas, dc.nama_fakultas) AS NamaFakultas,
			COALESCE(ul.kode_prodi, dc.kode_prodi) AS KodeProdi,
			COALESCE(ul.nama_prodi, dc.nama_prodi) AS NamaProdi,
			COALESCE(ul.role, dc.role) AS Role
		`).
		Joins(`
			LEFT JOIN (?) ul
				ON ul.id = CAST(e.createdByRef AS CHAR)
				AND LOWER(e.createdBy) = 'local'
		`, accountSub).
		Joins(`
			LEFT JOIN (?) dc
				ON dc.nidn = CAST(e.createdByRef AS CHAR)
				AND LOWER(e.createdBy) = 'simak'
		`, dosenSub).
		Where("e.id_bank_soal = ?", row.Id).
		Order("e.id DESC").
		Find(&extRows).Error

	if err != nil {
		return nil, err
	}

	row.ListExt = append(
		[]domainbanksoal.BankSoalExtDefault{},
		extRows...,
	)

	// =========================
	// PARSE TARGET UUID
	// =========================
	if row.RawTargetUUIDs != "" {
		uuids := strings.Split(row.RawTargetUUIDs, ",")

		for _, s := range uuids {
			u, err := uuid.Parse(strings.TrimSpace(s))
			if err == nil {
				row.TargetPertanyaan = append(row.TargetPertanyaan, u)
			}
		}
	}

	return &row, nil
}

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	"judul":         "b.judul",
	"semester":      "b.semester",
	"status":        "b.status",
	"role":          "COALESCE(ul.role, dc.role)",
	"kode_fakultas": "COALESCE(ul.kode_fakultas, dc.kode_fakultas)",
	"kode_prodi":    "COALESCE(ul.kode_prodi, dc.kode_prodi)",
	"nama_prodi": `
		CONCAT(
			COALESCE(ul.nama_prodi, dc.nama_prodi),
			' ',
			CASE COALESCE(ul.kode_jenjang, dc.kode_jenjang)
				WHEN 'E' THEN 'D3'
				WHEN 'A' THEN 'S3'
				WHEN 'B' THEN 'S2'
				WHEN 'C' THEN 'S1'
				ELSE ''
			END
		)
	`,
	"nama_fakultas": "COALESCE(ul.nama_fakultas, dc.nama_fakultas)",
}

func (r *BankSoalRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomain.SearchFilter,
	TargetFakultas string,
	TargetProdi string,
	TargetUnit string,
	TargetStatus string,
	page, limit *int,
	deleted bool,
	active bool,
) ([]domainbanksoal.BankSoalDefault, int64, error) {

	var rows = make([]domainbanksoal.BankSoalDefault, 0)
	var total int64

	// Subquery dosen
	dosenSub := r.db.
		Table("m_dosen d").
		Select(`
			CAST(d.NIDN AS CHAR) AS nidn,
			d.nama_dosen,
			f.kode_fakultas,
			f.nama_fakultas,
			p.kode_prodi,
			p.kode_jenjang,
			CONCAT(
				p.nama_prodi,
				CASE p.kode_jenjang
					WHEN 'J' THEN ' (Profesi)'
					WHEN 'E' THEN ' (D3)'
					WHEN 'D' THEN ' (D4)'
					WHEN 'A' THEN ' (S3)'
					WHEN 'B' THEN ' (S2)'
					WHEN 'C' THEN ' (S1)'
					ELSE ''
				END
			) AS nama_prodi,
			'dosen' as role
		`).
		Joins("LEFT JOIN m_fakultas f ON d.kode_fak = f.kode_fakultas").
		Joins("LEFT JOIN m_program_studi p ON d.kode_prodi = p.kode_prodi")

	// Subquery account
	accountSub := r.db.
		Table("users u").
		Select(`
			CAST(u.id AS CHAR) AS id,
			u.name,
			f.kode_fakultas,
			f.nama_fakultas,
			p.kode_prodi,
			p.kode_jenjang,
			CONCAT(
				p.nama_prodi,
				CASE p.kode_jenjang
					WHEN 'J' THEN ' (Profesi)'
					WHEN 'E' THEN ' (D3)'
					WHEN 'D' THEN ' (D4)'
					WHEN 'A' THEN ' (S3)'
					WHEN 'B' THEN ' (S2)'
					WHEN 'C' THEN ' (S1)'
					ELSE ''
				END
			) AS nama_prodi,
			u.level as role
		`).
		Joins("LEFT JOIN m_fakultas f ON u.fakultas = f.kode_fakultas").
		Joins("LEFT JOIN m_program_studi p ON u.prodi = p.kode_prodi")

	// Subquery pertanyaan
	pertanyaanSub := r.db.
		Table("template_pertanyaanv2 tp").
		Select(`
			tp.id_bank_soal,
			COUNT(tp.id) as total_pertanyaan,
			GROUP_CONCAT(tp.uuid) as uuids
		`).
		Where("tp.deleted_at IS NULL")

	// Filter status
	if TargetStatus != "" {
		pertanyaanSub = pertanyaanSub.Where("tp.status = ?", TargetStatus)
	} else {
		pertanyaanSub = pertanyaanSub.Where("tp.status = 'active'")
	}

	// Filter fakultas/prodi
	if TargetFakultas != "" && TargetProdi != "" {
		pertanyaanSub = pertanyaanSub.Where(`
			(tp.fakultas IS NULL OR tp.fakultas = ?)
			AND (tp.prodi IS NULL OR tp.prodi = ?)
		`, TargetFakultas, TargetProdi)
	} else if TargetFakultas != "" {
		pertanyaanSub = pertanyaanSub.Where(`
			tp.fakultas IS NULL OR tp.fakultas = ?
		`, TargetFakultas)
	} else if TargetUnit != "" {
		pertanyaanSub = pertanyaanSub.Where(`
			tp.unit IS NULL OR tp.unit = ?
		`, TargetUnit)
	}

	pertanyaanSub = pertanyaanSub.Group("tp.id_bank_soal")

	// Main query
	db := r.db.WithContext(ctx).
		Table("bank_soalv2 b").
		Debug().
		Select(`
			b.id as ID,
			b.uuid as UUID,
			b.judul as Judul,
			b.content as Content,
			b.deskripsi as Deskripsi,
			b.semester as Semester,
			b.tanggal_mulai as TanggalMulai,
			b.tanggal_akhir as TanggalAkhir,
			b.createdBy as CreatedBy,
			b.createdByRef AS CreatedByRef,
			b.deleted_at as DeletedAt,
			b.status as Status,

			COALESCE(ul.name, dc.nama_dosen) AS Nama,
			COALESCE(ul.role, dc.role) AS Role,
			COALESCE(ul.kode_fakultas, dc.kode_fakultas) AS KodeFakultas,
			COALESCE(ul.nama_fakultas, dc.nama_fakultas) AS NamaFakultas,
			COALESCE(ul.kode_prodi, dc.kode_prodi) AS KodeProdi,
			COALESCE(ul.nama_prodi, dc.nama_prodi) AS NamaProdi,

			COALESCE(pc.total_pertanyaan, 0) AS TotalPertanyaan,
			0 AS TotalInput,
			COALESCE(pc.uuids, '') AS TargetPertanyaan
		`).
		Joins(`LEFT JOIN (?) ul ON ul.id = CAST(b.createdByRef AS CHAR) AND LOWER(b.createdBy) = 'local'`, accountSub).
		Joins(`LEFT JOIN (?) dc ON dc.nidn = CAST(b.createdByRef AS CHAR) AND LOWER(b.createdBy) = 'simak'`, dosenSub).
		Joins(`LEFT JOIN (?) pc ON b.id = pc.id_bank_soal`, pertanyaanSub)

	// Filter active / deleted
	if active {
		db = db.Where("COALESCE(pc.total_pertanyaan,0) > 0")
	}
	if deleted {
		db = db.Where("b.deleted_at IS NOT NULL")
	} else {
		db = db.Where("b.deleted_at IS NULL")
	}

	// Advanced search filters
	for _, f := range searchFilters {
		col, ok := allowedSearchColumns[strings.ToLower(f.Field)]
		if !ok || f.Value == nil {
			continue
		}
		val := strings.TrimSpace(*f.Value)
		if val == "" {
			continue
		}
		switch strings.ToLower(f.Operator) {
		case "eq":
			db = db.Where(fmt.Sprintf("%s = ?", col), val)
		case "neq":
			db = db.Where(fmt.Sprintf("%s != ?", col), val)
		case "like":
			db = db.Where(fmt.Sprintf("%s LIKE ?", col), "%"+helper.EscapeLike(val)+"%")
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
				db = db.Where(fmt.Sprintf("%s IN ?", col), vals)
			}
		}
	}

	// Global search
	if s := strings.TrimSpace(search); s != "" {
		like := "%" + s + "%"
		var conditions []string
		var args []interface{}
		for _, col := range allowedSearchColumns {
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", col))
			args = append(args, like)
		}
		db = db.Where("("+strings.Join(conditions, " OR ")+")", args...)
	}

	// Hitung total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	db = db.Order("b.id DESC")
	if page != nil && limit != nil && *limit > 0 {
		offset := (*page - 1) * (*limit)
		db = db.Offset(offset).Limit(*limit)
	}

	// Fetch rows
	if err := db.Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	// =========================
	// LOAD EXTENSION DATA
	// =========================
	if len(rows) > 0 {

		ids := make([]uint, 0, len(rows))
		for _, row := range rows {
			ids = append(ids, row.Id)
		}

		extRows := make([]domainbanksoal.BankSoalExtDefault, 0, len(rows))

		err := r.db.
			Table("bank_soal_extendv2 e").
			Debug().
			Select(`
				e.id as ID,
				e.uuid as UUID,
				e.id_bank_soal as IdBankSoal,
				e.tanggal_mulai as TanggalMulai,
				e.tanggal_akhir as TanggalAkhir,
				e.createdBy as CreatedBy,
				e.createdByRef as CreatedByRef,

				COALESCE(ul.kode_fakultas, dc.kode_fakultas) AS KodeFakultas,
				COALESCE(ul.nama_fakultas, dc.nama_fakultas) AS NamaFakultas,
				COALESCE(ul.kode_prodi, dc.kode_prodi) AS KodeProdi,
				COALESCE(ul.nama_prodi, dc.nama_prodi) AS NamaProdi,
				COALESCE(ul.role, dc.role) AS Role
			`).
			Joins(`LEFT JOIN (?) ul ON ul.id = CAST(e.createdByRef AS CHAR) AND LOWER(e.createdBy) = 'local'`, accountSub).
			Joins(`LEFT JOIN (?) dc ON dc.nidn = CAST(e.createdByRef AS CHAR) AND LOWER(e.createdBy) = 'simak'`, dosenSub).
			Where("e.id_bank_soal IN ?", ids).
			Order("e.id DESC").
			Find(&extRows).Error

		if err != nil {
			return nil, 0, err
		}

		extMap := make(map[uint][]domainbanksoal.BankSoalExtDefault, 0)

		for _, item := range extRows {
			extMap[item.IdBankSoal] = append(
				extMap[item.IdBankSoal],
				item,
			)
		}

		for i := range rows {
			if val, ok := extMap[rows[i].Id]; ok {
				rows[i].ListExt = append(
					[]domainbanksoal.BankSoalExtDefault{},
					val...,
				)
			} else {
				rows[i].ListExt = []domainbanksoal.BankSoalExtDefault{}
			}
		}
	}

	// Parse UUID TargetPertanyaan
	for i := range rows {
		if rows[i].RawTargetUUIDs != "" {
			uuids := strings.SplitSeq(rows[i].RawTargetUUIDs, ",")
			for s := range uuids {
				if u, err := uuid.Parse(strings.TrimSpace(s)); err == nil {
					rows[i].TargetPertanyaan = append(rows[i].TargetPertanyaan, u)
				}
			}
		}
	}

	return rows, total, nil
}

// ------------------------
// CREATE
// ------------------------
func (r *BankSoalRepository) Create(ctx context.Context, banksoal *domainbanksoal.BankSoal) error {
	return r.db.WithContext(ctx).Create(banksoal).Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *BankSoalRepository) Update(ctx context.Context, banksoal *domainbanksoal.BankSoal) error {
	return r.db.WithContext(ctx).Save(banksoal).Error
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *BankSoalRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domainbanksoal.BankSoal{}).Error
}

func (r *BankSoalRepository) SetupUuid(ctx context.Context) error {
	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domainbanksoal.BankSoal{}).
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
			"UPDATE bank_soalv2 SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *BankSoalRepository) CountCopy(ctx context.Context, judul string) (int, error) {

	var count int64

	err := r.db.WithContext(ctx).
		Table("bank_soalv2").
		Where("judul = ? OR judul LIKE ?",
			fmt.Sprintf("salin - %s", judul),
			fmt.Sprintf("salin (%%) - %s", judul),
		).
		Count(&count).Error

	return int(count), err
}

// ------------------------
// CREATE EXT
// ------------------------
func (r *BankSoalRepository) CreateExt(ctx context.Context, banksoalext *domainbanksoal.BankSoalExt) error {
	return r.db.WithContext(ctx).Create(banksoalext).Error
}

// ------------------------
// DELETE EXT (by UUID)
// ------------------------
func (r *BankSoalRepository) DeleteExt(ctx context.Context, uid uuid.UUID, idbanksoal uint) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domainbanksoal.BankSoalExt{}).Error
}
