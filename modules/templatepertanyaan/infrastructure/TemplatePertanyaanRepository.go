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
	db      *gorm.DB
	dbSimak *gorm.DB
}

func NewTemplatePertanyaanRepository(db *gorm.DB, dbSimak ...*gorm.DB) domaintemplatepertanyaan.ITemplatePertanyaanRepository {
	var simak *gorm.DB
	if len(dbSimak) > 0 {
		simak = dbSimak[0]
	}
	return &TemplatePertanyaanRepository{db: db, dbSimak: simak}
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
		Joins("LEFT JOIN users u ON a.createdByRef = u.id").
		Joins("LEFT JOIN m_fakultas f ON f.kode_fakultas = u.fakultas").
		Joins("LEFT JOIN m_program_studi p ON p.kode_prodi = u.prodi").
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
			CASE
				WHEN u.prodi IS NOT NULL AND u.prodi != '' THEN CONCAT("PRODI ", u.prodi)
				WHEN u.fakultas IS NOT NULL AND u.fakultas != '' THEN CONCAT("FAKULTAS ", u.fakultas)
				WHEN LOWER(u.level) = 'fakultas' THEN 'fakultas'
				WHEN LOWER(u.level) = 'prodi' THEN 'prodi'
				ELSE COALESCE(NULLIF(TRIM(a.createdBy), ''), 'admin')
			END as CreatedBy,
			a.createdByRef as CreatedByRef,
			COALESCE(NULLIF(TRIM(a.fakultas), ''), f.nama_fakultas, f.kode_fakultas, u.fakultas) as Fakultas,
			COALESCE(NULLIF(TRIM(a.prodi), ''), p.nama_prodi, p.kode_prodi, u.prodi) as Prodi,
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
		Joins("LEFT JOIN users u ON a.createdByRef = u.id").
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
			CASE
				WHEN u.prodi IS NOT NULL AND u.prodi != '' THEN CONCAT("PRODI ", u.prodi)
				WHEN u.fakultas IS NOT NULL AND u.fakultas != '' THEN CONCAT("FAKULTAS ", u.fakultas)
				WHEN LOWER(u.level) = 'fakultas' THEN 'fakultas'
				WHEN LOWER(u.level) = 'prodi' THEN 'prodi'
				ELSE COALESCE(NULLIF(TRIM(a.createdBy), ''), 'admin')
			END as CreatedBy,
			a.createdByRef as CreatedByRef,
			COALESCE(NULLIF(TRIM(a.fakultas), ''), u.fakultas) as Fakultas,
			COALESCE(NULLIF(TRIM(a.prodi), ''), u.prodi) as Prodi,
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
			WHEN u.prodi IS NOT NULL AND u.prodi != '' THEN CONCAT("PRODI ", u.prodi)
			WHEN u.fakultas IS NOT NULL AND u.fakultas != '' THEN CONCAT("FAKULTAS ", u.fakultas)
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

	r.resolveFakultasProdiNames(ctx, rows)

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
// DELETE (SOFT DELETE)
// ------------------------
func (r *TemplatePertanyaanRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domaintemplatepertanyaan.TemplatePertanyaan{}).Error
}

// ------------------------
// RESTORE
// ------------------------
func (r *TemplatePertanyaanRepository) Restore(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&domaintemplatepertanyaan.TemplatePertanyaan{}).
		Unscoped().
		Where("uuid = ?", uid).
		Update("deleted_at", nil).Error
}

// ------------------------
// SETUP UUID
// ------------------------
func (r *TemplatePertanyaanRepository) SetupUuid(ctx context.Context) error {
	var list []domaintemplatepertanyaan.TemplatePertanyaan
	if err := r.db.WithContext(ctx).Where("uuid IS NULL OR uuid = ''").Find(&list).Error; err != nil {
		return err
	}

	for _, item := range list {
		item.UUID = uuid.New()
		if err := r.db.WithContext(ctx).Save(&item).Error; err != nil {
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

// ------------------------
// COPY TEMPLATE PERTANYAAN
// ------------------------
func (r *TemplatePertanyaanRepository) CopyByBankSoal(
	ctx context.Context,
	tx *gorm.DB,
	sourceBankSoalID uint,
	targetBankSoalID uint,
	resource string,
	sid string,
	isDefault ...bool,
) (map[uint]uint, error) {

	checkDefault := false
	if len(isDefault) > 0 {
		checkDefault = isDefault[0]
	}

	var rows []domaintemplatepertanyaan.TemplatePertanyaan

	if err := tx.WithContext(ctx).
		Where("id_bank_soal = ? AND deleted_at IS NULL", sourceBankSoalID).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	mapping := make(map[uint]uint)

	for _, row := range rows {

		oldID := row.ID

		row.ID = 0
		row.UUID = uuid.New()

		row.IdBankSoal = targetBankSoalID
		if !checkDefault {
			row.CreatedBy = helper.StrPtr(resource)
			row.CreatedByRef = helper.StrPtr(sid)
		}

		if err := tx.WithContext(ctx).
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

func (r *TemplatePertanyaanRepository) resolveFakultasProdiNames(ctx context.Context, rows []domaintemplatepertanyaan.TemplatePertanyaanDefault) {
	if r.dbSimak == nil || len(rows) == 0 {
		return
	}

	fakCodesMap := make(map[string]bool)
	prodiCodesMap := make(map[string]bool)

	for _, row := range rows {
		if row.Fakultas != nil && *row.Fakultas != "" {
			for _, code := range strings.Split(*row.Fakultas, ",") {
				if c := strings.TrimSpace(code); c != "" {
					fakCodesMap[c] = true
				}
			}
		}
		if row.Prodi != nil && *row.Prodi != "" {
			for _, code := range strings.Split(*row.Prodi, ",") {
				if c := strings.TrimSpace(code); c != "" {
					prodiCodesMap[c] = true
				}
			}
		}

		if row.CreatedBy != nil {
			createdByVal := *row.CreatedBy
			if strings.HasPrefix(createdByVal, "FAKULTAS ") {
				c := strings.TrimSpace(strings.TrimPrefix(createdByVal, "FAKULTAS "))
				if c != "" {
					fakCodesMap[c] = true
				}
			} else if strings.HasPrefix(createdByVal, "PRODI ") {
				c := strings.TrimSpace(strings.TrimPrefix(createdByVal, "PRODI "))
				if c != "" {
					prodiCodesMap[c] = true
				}
			}
		}
	}

	fakNameMap := make(map[string]string)
	if len(fakCodesMap) > 0 {
		var fakList []struct {
			KodeFakultas string `gorm:"column:kode_fakultas"`
			NamaFakultas string `gorm:"column:nama_fakultas"`
		}
		codes := make([]string, 0, len(fakCodesMap))
		for c := range fakCodesMap {
			codes = append(codes, c)
		}
		_ = r.dbSimak.WithContext(ctx).
			Table("m_fakultas").
			Select("kode_fakultas, nama_fakultas").
			Where("kode_fakultas IN ?", codes).
			Find(&fakList).Error

		for _, f := range fakList {
			fakNameMap[strings.TrimSpace(f.KodeFakultas)] = strings.TrimSpace(f.NamaFakultas)
		}
	}

	prodiNameMap := make(map[string]string)
	if len(prodiCodesMap) > 0 {
		var prodiList []struct {
			KodeProdi string `gorm:"column:kode_prodi"`
			NamaProdi string `gorm:"column:nama_prodi"`
			Jenjang   string `gorm:"column:kode_jenjang"`
		}
		codes := make([]string, 0, len(prodiCodesMap))
		for c := range prodiCodesMap {
			codes = append(codes, c)
		}
		_ = r.dbSimak.WithContext(ctx).
			Table("m_program_studi").
			Select("kode_prodi, nama_prodi, kode_jenjang").
			Where("kode_prodi IN ?", codes).
			Find(&prodiList).Error

		for _, p := range prodiList {
			prodiNameMap[strings.TrimSpace(p.KodeProdi)] = strings.TrimSpace(p.NamaProdi)
		}
	}

	for i := range rows {
		resolvedFak := ""
		resolvedProdi := ""

		if rows[i].Fakultas != nil && *rows[i].Fakultas != "" {
			parts := strings.Split(*rows[i].Fakultas, ",")
			var names []string
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if name, ok := fakNameMap[p]; ok && name != "" {
					names = append(names, name)
				} else if p != "" {
					names = append(names, p)
				}
			}
			resolvedFak = strings.Join(names, ", ")
			rows[i].Fakultas = &resolvedFak
		}

		if rows[i].Prodi != nil && *rows[i].Prodi != "" {
			parts := strings.Split(*rows[i].Prodi, ",")
			var names []string
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if name, ok := prodiNameMap[p]; ok && name != "" {
					names = append(names, name)
				} else if p != "" {
					names = append(names, p)
				}
			}
			resolvedProdi = strings.Join(names, ", ")
			rows[i].Prodi = &resolvedProdi
		}

		if rows[i].CreatedBy != nil {
			createdByVal := *rows[i].CreatedBy
			if strings.HasPrefix(createdByVal, "FAKULTAS ") {
				codeVal := strings.TrimSpace(strings.TrimPrefix(createdByVal, "FAKULTAS "))
				if name, ok := fakNameMap[codeVal]; ok && name != "" {
					val := "FAKULTAS " + name
					rows[i].CreatedBy = &val
				} else if resolvedFak != "" {
					val := "FAKULTAS " + resolvedFak
					rows[i].CreatedBy = &val
				}
			} else if strings.HasPrefix(createdByVal, "PRODI ") {
				codeVal := strings.TrimSpace(strings.TrimPrefix(createdByVal, "PRODI "))
				if name, ok := prodiNameMap[codeVal]; ok && name != "" {
					val := "PRODI " + name
					rows[i].CreatedBy = &val
				} else if resolvedProdi != "" {
					val := "PRODI " + resolvedProdi
					rows[i].CreatedBy = &val
				} else {
					val := "PRODI "
					rows[i].CreatedBy = &val
				}
			}
		}
	}
}
