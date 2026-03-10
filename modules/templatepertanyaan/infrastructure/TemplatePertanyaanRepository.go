package infrastructure

import (
	commondomaintemplatepertanyaan "UnpakSiamida/common/domain"
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
		Joins("LEFT JOIN kategori k ON k.id = a.id_kategori").
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
		a.required as Required,
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

	return &rowData, nil
}

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	"pertanyaan": "a.pertanyaan",
	"kategori":   "k.nama_kategori",
}

// ------------------------
// GET ALL
// ------------------------
func (r *TemplatePertanyaanRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomaintemplatepertanyaan.SearchFilter,
	page, limit *int,
	deleted bool,
) ([]domaintemplatepertanyaan.TemplatePertanyaanDefault, int64, error) {

	var rows = make([]domaintemplatepertanyaan.TemplatePertanyaanDefault, 0)
	var total int64

	db := r.db.Debug().WithContext(ctx).
		Table("template_pertanyaanv2 a").
		Joins("LEFT JOIN kategori k ON k.id = a.id_kategori").
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
		b.uuid as UuidKategori,
		k.nama_kategori as Kategori,
		a.required as Required,
		a.created_at as CreatedAt,
		a.updated_at as UpdatedAt
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
