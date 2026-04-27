package infrastructure

import (
	commondomainkategori "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	domainkategori "UnpakSiamida/modules/kategori/domain"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type KategoriRepository struct {
	db *gorm.DB
}

func NewKategoriRepository(db *gorm.DB) domainkategori.IKategoriRepository {
	return &KategoriRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *KategoriRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
	var Kategori domainkategori.Kategori

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&Kategori).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &Kategori, nil
}

// ------------------------
// GET DEFAULT BY UUID
// ------------------------
func (r *KategoriRepository) GetDefaultByUuid(
	ctx context.Context,
	id uuid.UUID,
) (*domainkategori.KategoriDefault, error) {

	// Ambil hanya kolom yang benar-benar ada di struct KategoriDefault
	var rowData domainkategori.KategoriDefault

	uuidSub := clause.Expr{
		SQL: `(SELECT k.uuid FROM kategori k WHERE k.id = a.sub_kategori LIMIT 1)`,
	}

	namaSub := clause.Expr{
		SQL: `(SELECT k.nama_kategori FROM kategori k WHERE k.id = a.sub_kategori LIMIT 1)`,
	}

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
					WHEN 'E' THEN ' (D3)'
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
					WHEN 'E' THEN ' (D3)'
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

	err := r.db.WithContext(ctx).
		Table("kategori a").
		Select(`
			a.id as ID,
			a.uuid as UUID,
			a.nama_kategori as NamaKategori,
			a.full_text as FullTexts,
			a.sub_kategori as IdSubKategori,
			? as UuidSubKategori,
			? as NamaSubKategori,
			a.createdBy as CreatedBy,
			a.createdByRef as CreatedByRef,
			a.deleted_at as DeletedAt,

			COALESCE(ul.name, dc.nama_dosen) AS Nama,
			COALESCE(ul.role, dc.role) AS Role,
			COALESCE(ul.kode_fakultas, dc.kode_fakultas) AS KodeFakultas,
			COALESCE(ul.nama_fakultas, dc.nama_fakultas) AS NamaFakultas,
			COALESCE(ul.kode_prodi, dc.kode_prodi) AS KodeProdi,
			COALESCE(ul.nama_prodi, dc.nama_prodi) AS NamaProdi
	`, uuidSub, namaSub).
		Joins(`LEFT JOIN (?) ul ON ul.id = CAST(a.createdByRef AS CHAR) AND LOWER(a.createdBy) = 'local'`, accountSub).
		Joins(`LEFT JOIN (?) dc ON dc.nidn = CAST(a.createdByRef AS CHAR) AND LOWER(a.createdBy) = 'simak'`, dosenSub).
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

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	"kategori":  "a.nama_kategori",
	"full_text": "a.full_text",

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

// ------------------------
// GET ALL
// ------------------------
func (r *KategoriRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomainkategori.SearchFilter,
	page, limit *int,
	deleted bool,
) ([]domainkategori.KategoriDefault, int64, error) {

	var rows = make([]domainkategori.KategoriDefault, 0)
	var total int64

	uuidSub := clause.Expr{
		SQL: `(SELECT k.uuid FROM kategori k WHERE k.id = a.sub_kategori LIMIT 1)`,
	}

	namaSub := clause.Expr{
		SQL: `(SELECT k.nama_kategori FROM kategori k WHERE k.id = a.sub_kategori LIMIT 1)`,
	}

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
					WHEN 'E' THEN ' (D3)'
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
					WHEN 'E' THEN ' (D3)'
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

	db := r.db.Debug().WithContext(ctx).
		Table("kategori a").
		Select(`
			a.id as ID,
			a.uuid as UUID,
			a.nama_kategori as NamaKategori,
			a.full_text as FullTexts,
			a.sub_kategori as IdSubKategori,
			? as UuidSubKategori,
			? as NamaSubKategori,
			a.createdBy as CreatedBy,
			a.createdByRef as CreatedByRef,
			a.deleted_at as DeletedAt,

			COALESCE(ul.name, dc.nama_dosen) AS Nama,
			COALESCE(ul.role, dc.role) AS Role,
			COALESCE(ul.kode_fakultas, dc.kode_fakultas) AS KodeFakultas,
			COALESCE(ul.nama_fakultas, dc.nama_fakultas) AS NamaFakultas,
			COALESCE(ul.kode_prodi, dc.kode_prodi) AS KodeProdi,
			COALESCE(ul.nama_prodi, dc.nama_prodi) AS NamaProdi
	`, uuidSub, namaSub).
		Joins(`LEFT JOIN (?) ul ON ul.id = CAST(a.createdByRef AS CHAR) AND LOWER(a.createdBy) = 'local'`, accountSub).
		Joins(`LEFT JOIN (?) dc ON dc.nidn = CAST(a.createdByRef AS CHAR) AND LOWER(a.createdBy) = 'simak'`, dosenSub)

	if deleted {
		db = db.Where(clause.Expr{
			SQL: "a.deleted_at IS NOT NULL",
		})
	} else {
		db = db.Where(clause.Expr{
			SQL: "a.deleted_at IS NULL",
		})
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
func (r *KategoriRepository) Create(ctx context.Context, kategori *domainkategori.Kategori) error {
	return r.db.WithContext(ctx).Create(kategori).Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *KategoriRepository) Update(ctx context.Context, kategori *domainkategori.Kategori) error {
	return r.db.WithContext(ctx).Save(kategori).Error
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *KategoriRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domainkategori.Kategori{}).Error
}

func (r *KategoriRepository) SetupUuid(ctx context.Context) error {
	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domainkategori.Kategori{}).
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
			"UPDATE kategori SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *KategoriRepository) CountCopy(ctx context.Context, nama_kategori string) (int, error) {

	var count int64

	err := r.db.WithContext(ctx).
		Table("kategori").
		Where("nama_kategori = ? OR nama_kategori LIKE ?",
			fmt.Sprintf("salin - %s", nama_kategori),
			fmt.Sprintf("salin (%%) - %s", nama_kategori),
		).
		Count(&count).Error

	return int(count), err
}

func (r *KategoriRepository) GetChildren(ctx context.Context, parentID int) ([]domainkategori.Kategori, error) {

	rows := make([]domainkategori.Kategori, 0)

	err := r.db.WithContext(ctx).
		Where("sub_kategori = ?", parentID).
		Find(&rows).Error

	return rows, err
}

func (r *KategoriRepository) RebuildFullText(ctx context.Context) error {

	list := make([]domainkategori.Kategori, 0)

	if err := r.db.WithContext(ctx).
		Find(&list).Error; err != nil {
		return err
	}

	childrenMap := map[uint][]domainkategori.Kategori{}
	roots := make([]domainkategori.Kategori, 0)

	for _, k := range list {

		if k.SubKategori == nil {
			roots = append(roots, k)
			continue
		}

		childrenMap[*k.SubKategori] = append(childrenMap[*k.SubKategori], k)
	}

	for _, root := range roots {

		if err := r.buildNodeFast(
			ctx,
			root,
			root.NamaKategori,
			childrenMap,
		); err != nil {
			return err
		}
	}

	return nil
}

func (r *KategoriRepository) buildNodeFast(
	ctx context.Context,
	node domainkategori.Kategori,
	path string,
	childrenMap map[uint][]domainkategori.Kategori,
) error {

	if err := r.db.WithContext(ctx).
		Model(&domainkategori.Kategori{}).
		Where("id = ?", node.ID).
		Update("full_text", path).Error; err != nil {
		return err
	}

	children := childrenMap[node.ID]

	for _, child := range children {

		newPath := fmt.Sprintf("%s > %s", path, child.NamaKategori)

		if err := r.buildNodeFast(
			ctx,
			child,
			newPath,
			childrenMap,
		); err != nil {
			return err
		}
	}

	return nil
}

func (r *KategoriRepository) WithTx(
	ctx context.Context,
	fn func(repo domainkategori.IKategoriRepository) error,
) error {

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		txRepo := &KategoriRepository{
			db: tx,
		}

		return fn(txRepo)
	})
}
