package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	domaintemplatepertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TemplatePertanyaanRepository struct {
	db *gorm.DB
}

func NewTemplatePertanyaanRepository(db *gorm.DB) domaintemplatepertanyaan.ITemplatePertanyaanRepository {
	return &TemplatePertanyaanRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *TemplatePertanyaanRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domaintemplatepertanyaan.TemplatePertanyaan, error) {
	var TemplatePertanyaan domaintemplatepertanyaan.TemplatePertanyaan

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&TemplatePertanyaan).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &TemplatePertanyaan, nil
}

// ------------------------
// GET DEFAULT BY UUID
// ------------------------
func (r *TemplatePertanyaanRepository) GetDefaultByUuid(
	ctx context.Context,
	id uuid.UUID,
) (*domaintemplatepertanyaan.TemplatePertanyaanDefault, error) {

	// Ambil hanya kolom yang benar-benar ada di struct TemplatePertanyaanDefault
	var rowData domaintemplatepertanyaan.TemplatePertanyaanDefault

	err := r.db.Debug().WithContext(ctx).
		Table("template_pertanyaanv2 a").
		Joins("LEFT JOIN kategoriv2 k ON k.id = a.id_kategori").
		Joins("LEFT JOIN bank_soalv2 b ON b.id = a.id_bank_soal").
		Select(`
		a.id as ID,
		a.uuid as UUID,
		a.id_bank_soal as IdBankSoal,
		b.uuid as UuidBankSoal,
		b.judul as NamaBankSoal,
		a.pertanyaan as Pertanyaan,
		a.jenis_pilihan as JenisPilihan,
		a.bobot as Bobot,
		a.id_kategori as IdKategori,
		k.uuid as UuidKategori,
		k.nama_kategori as Kategori,
		k.full_text as FullPath,
		a.status as Status,
		a.required as Required,
		a.created_at as CreatedAt,
		a.updated_at as UpdatedAt,
		a.deleted_at as DeletedAt
	`).
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

func (r *TemplatePertanyaanRepository) GetDefaultWithAnswareByUuid(
	ctx context.Context,
	id uuid.UUID,
) (*domaintemplatepertanyaan.TemplatePertanyaanWithAnswareDefault, error) {

	var rowData domaintemplatepertanyaan.TemplatePertanyaanWithAnswareDefault

	// =========================
	// QUERY HEADER PERTANYAAN
	// =========================
	err := r.db.Debug().
		WithContext(ctx).
		Table("template_pertanyaanv2 a").
		Joins("LEFT JOIN kategoriv2 k ON k.id = a.id_kategori").
		Joins("LEFT JOIN bank_soalv2 b ON b.id = a.id_bank_soal").
		Select(`
			a.id as ID,
			a.uuid as UUID,
			a.id_bank_soal as IdBankSoal,
			b.uuid as UUIDBankSoal,
			b.judul as NamaBankSoal,
			a.pertanyaan as Pertanyaan,
			a.jenis_pilihan as JenisPilihan,
			a.bobot as Bobot,
			a.id_kategori as IdKategori,
			k.uuid as UUIDKategori,
			k.nama_kategori as Kategori,
			k.full_text as FullPath,
			a.required as Required,
			a.status as Status,
			a.createdBy as CreatedBy,
			a.createdByRef as CreatedByRef,
			a.fakultas as Fakultas,
			a.prodi as Prodi,
			a.unit as Unit,
			a.jenjang as Jenjang,
			a.deleted_at as DeletedAt,
			a.created_at as CreatedAt,
			a.updated_at as UpdatedAt
		`).
		Where("a.uuid = ?", id).
		Take(&rowData).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	// =========================
	// QUERY LIST JAWABAN
	// =========================
	listJawaban := make([]domaintemplatepertanyaan.TemplateJawabanDefault, 0)

	err = r.db.Debug().
		WithContext(ctx).
		Table("template_pilihanv2 p").
		Select(`
			p.id as ID,
			p.uuid as UUID,
			p.id_template_pertanyaan as IdTemplatePertanyaan,
			? as UUIDTemplatePertanyaan,
			? as NamaTemplatePertanyaan,
			p.jawaban as Jawaban,
			p.nilai as Nilai,
			p.isFreeText as IsFreeText,
			p.deleted_at as DeletedAt,
			p.created_at as CreatedAt,
			p.updated_at as UpdatedAt
		`, rowData.UUID, rowData.Pertanyaan).
		Where("p.id_template_pertanyaan = ?", rowData.ID).
		Where("p.deleted_at IS NULL").
		Order("p.id ASC").
		Find(&listJawaban).Error

	if err != nil {
		return nil, err
	}

	rowData.ListJawaban = listJawaban

	return &rowData, nil
}

func (r *TemplatePertanyaanRepository) GetDefaultWithAnswareByBankSoal(
	ctx context.Context,
	id_bank_soal uint,
) ([]domaintemplatepertanyaan.TemplatePertanyaanWithAnswareDefault, error) {

	rows := make([]domaintemplatepertanyaan.TemplatePertanyaanWithAnswareDefault, 0)

	// =========================
	// QUERY LIST PERTANYAAN
	// =========================
	err := r.db.Debug().
		WithContext(ctx).
		Table("template_pertanyaanv2 a").
		Joins("LEFT JOIN kategoriv2 k ON k.id = a.id_kategori").
		Joins("LEFT JOIN bank_soalv2 b ON b.id = a.id_bank_soal").
		Select(`
			a.id as ID,
			a.uuid as UUID,
			a.id_bank_soal as IdBankSoal,
			b.uuid as UUIDBankSoal,
			b.judul as NamaBankSoal,
			a.pertanyaan as Pertanyaan,
			a.jenis_pilihan as JenisPilihan,
			a.bobot as Bobot,
			a.id_kategori as IdKategori,
			k.uuid as UUIDKategori,
			k.nama_kategori as Kategori,
			k.full_text as FullPath,
			a.required as Required,
			a.status as Status,
			a.createdBy as CreatedBy,
			a.createdByRef as CreatedByRef,
			a.fakultas as Fakultas,
			a.prodi as Prodi,
			a.unit as Unit,
			a.jenjang as Jenjang,
			a.deleted_at as DeletedAt,
			a.created_at as CreatedAt,
			a.updated_at as UpdatedAt
		`).
		Where("a.id_bank_soal = ?", id_bank_soal).
		Where("a.deleted_at IS NULL").
		Order("a.id ASC").
		Find(&rows).Error

	if err != nil {
		return nil, err
	}

	// =========================
	// PARALLEL FETCH JAWABAN
	// =========================
	type result struct {
		index   int
		jawaban []domaintemplatepertanyaan.TemplateJawabanDefault
		err     error
	}

	resultCh := make(chan result, len(rows))

	for i := range rows {
		go func(i int) {
			var listJawaban []domaintemplatepertanyaan.TemplateJawabanDefault

			err := r.db.Debug().
				WithContext(ctx).
				Table("template_pilihanv2 p").
				Select(`
					p.id as ID,
					p.uuid as UUID,
					p.id_template_pertanyaan as IdTemplatePertanyaan,
					? as UUIDTemplatePertanyaan,
					? as NamaTemplatePertanyaan,
					p.jawaban as Jawaban,
					p.nilai as Nilai,
					p.isFreeText as IsFreeText,
					p.deleted_at as DeletedAt,
					p.created_at as CreatedAt,
					p.updated_at as UpdatedAt
				`, rows[i].UUID, rows[i].Pertanyaan).
				Where("p.id_template_pertanyaan = ?", rows[i].ID).
				Where("p.deleted_at IS NULL").
				Order("p.id ASC").
				Find(&listJawaban).Error

			resultCh <- result{
				index:   i,
				jawaban: listJawaban,
				err:     err,
			}
		}(i)
	}

	// collect result
	for i := 0; i < len(rows); i++ {
		res := <-resultCh
		if res.err != nil {
			return nil, res.err
		}

		rows[res.index].ListJawaban = res.jawaban
	}

	return rows, nil
}

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	"uuidbanksoal": "b.uuid",
	"pertanyaan":   "a.pertanyaan",
	"kategori":     "k.nama_kategori",
}

// ------------------------
// GET ALL
// ------------------------
func (r *TemplatePertanyaanRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomain.SearchFilter,
	page, limit *int,
	deleted bool,
) ([]domaintemplatepertanyaan.TemplatePertanyaanDefault, int64, error) {

	var rows = make([]domaintemplatepertanyaan.TemplatePertanyaanDefault, 0)
	var total int64

	db := r.db.Debug().WithContext(ctx).
		Table("template_pertanyaanv2 a").
		Joins("LEFT JOIN kategoriv2 k ON k.id = a.id_kategori").
		Joins("LEFT JOIN bank_soalv2 b ON b.id = a.id_bank_soal").
		Joins("LEFT JOIN users u ON a.createdByRef = u.id").
		Joins("LEFT JOIN m_fakultas f ON u.fakultas = f.kode_fakultas").
		Joins("LEFT JOIN m_program_studi p ON u.prodi = p.kode_prodi").
		Select(`
		a.id as ID,
		a.uuid as UUID,
		a.id_bank_soal as IdBankSoal,
		b.uuid as UuidBankSoal,
		b.judul as NamaBankSoal,
		a.pertanyaan as Pertanyaan,
		a.jenis_pilihan as JenisPilihan,
		a.bobot as Bobot,
		a.id_kategori as IdKategori,
		k.uuid as UuidKategori,
		k.nama_kategori as Kategori,
		k.full_text as FullPath,
		a.status as Status,
		a.required as Required,
		a.fakultas as Fakultas,
		a.prodi as Prodi,
		a.unit as Unit,
		a.created_at as CreatedAt,
		CASE
			WHEN u.prodi IS NOT NULL OR u.prodi != '' THEN CONCAT(
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
			)
			WHEN u.fakultas IS NOT NULL OR u.fakultas != '' THEN CONCAT("FAKULTAS ",f.nama_fakultas)
			ELSE 'admin'
		END as CreatedBy,
		a.createdByRef as CreatedByRef,
		a.updated_at as UpdatedAt,
		a.deleted_at as DeletedAt
	`)

	if deleted {
		db = db.Where(clause.Expr{
			SQL: "a.deleted_at IS NOT NULL",
		})
	} else {
		db = db.Where(clause.Expr{
			SQL: "a.deleted_at IS NULL",
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

// ------------------------
// CREATE
// ------------------------
func (r *TemplatePertanyaanRepository) Create(ctx context.Context, templatepertanyaan *domaintemplatepertanyaan.TemplatePertanyaan) error {
	return r.db.WithContext(ctx).Create(templatepertanyaan).Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *TemplatePertanyaanRepository) Update(ctx context.Context, templatepertanyaan *domaintemplatepertanyaan.TemplatePertanyaan) error {
	return r.db.WithContext(ctx).Save(templatepertanyaan).Error
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *TemplatePertanyaanRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domaintemplatepertanyaan.TemplatePertanyaan{}).Error
}

func (r *TemplatePertanyaanRepository) SetupUuid(ctx context.Context) error {
	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domaintemplatepertanyaan.TemplatePertanyaan{}).
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
			"UPDATE template_pertanyaanv2 SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *TemplatePertanyaanRepository) CountCopy(ctx context.Context, judul string) (int, error) {

	var count int64

	err := r.db.WithContext(ctx).
		Table("template_pertanyaanv2").
		Where("pertanyaan = ? OR pertanyaan LIKE ?",
			fmt.Sprintf("salin - %s", judul),
			fmt.Sprintf("salin (%%) - %s", judul),
		).
		Count(&count).Error

	return int(count), err
}

func (r *TemplatePertanyaanRepository) CopyByBankSoal(
	ctx context.Context,
	tx *gorm.DB,
	sourceBankSoalID uint,
	targetBankSoalID uint,
	resource string,
	sid string,
) (map[uint]uint, error) {

	var rows []domaintemplatepertanyaan.TemplatePertanyaan

	if err := tx.WithContext(ctx).
		Where("id_bank_soal = ?", sourceBankSoalID).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return map[uint]uint{}, nil
	}

	mapping := make(map[uint]uint)

	for _, row := range rows {

		oldID := row.ID

		row.ID = 0
		row.UUID = uuid.New()

		row.IdBankSoal = targetBankSoalID
		row.CreatedBy = helper.StrPtr(resource)
		row.CreatedByRef = helper.StrPtr(sid)

		if err := r.db.WithContext(ctx).
			Create(&row).Error; err != nil {
			return nil, err
		}

		mapping[oldID] = row.ID
	}

	return mapping, nil
}

func (r *TemplatePertanyaanRepository) WithTx(tx any) domaintemplatepertanyaan.ITemplatePertanyaanRepository {
	return &TemplatePertanyaanRepository{
		db: tx.(*gorm.DB),
	}
}

func (r *TemplatePertanyaanRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	return r.db.WithContext(ctx).Begin(), nil
}
