package infrastructure

import (
	commondomaintemplatejawaban "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	domaintemplatejawaban "UnpakSiamida/modules/templatejawaban/domain"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TemplateJawabanRepository struct {
	db *gorm.DB
}

func NewTemplateJawabanRepository(db *gorm.DB) domaintemplatejawaban.ITemplateJawabanRepository {
	return &TemplateJawabanRepository{db: db}
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *TemplateJawabanRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
	var TemplateJawaban domaintemplatejawaban.TemplateJawaban

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		First(&TemplateJawaban).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }

	if err != nil {
		return nil, err
	}

	return &TemplateJawaban, nil
}

// ------------------------
// GET DEFAULT BY UUID
// ------------------------
func (r *TemplateJawabanRepository) GetDefaultByUuid(
	ctx context.Context,
	id uuid.UUID,
) (*domaintemplatejawaban.TemplateJawabanDefault, error) {

	// Ambil hanya kolom yang benar-benar ada di struct TemplateJawabanDefault
	var rowData domaintemplatejawaban.TemplateJawabanDefault

	err := r.db.Debug().WithContext(ctx).
		Table("template_pilihanv2 a").
		Joins("LEFT JOIN template_pertanyaanv2 b ON b.id = a.id_template_pertanyaan").
		Select(`
		a.id as ID,
		a.uuid as UUID,
		a.id_template_pertanyaan as IdTemplatePertanyaan,
		b.uuid as UUIDTemplatePertanyaan,
		b.pertanyaan as NamaTemplatePertanyaan,
		a.jawaban as Jawaban,
		a.nilai as Nilai,
		a.isFreeText as IsFreeText,
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

	return &rowData, nil
}

var allowedSearchColumns = map[string]string{
	// key:param -> db column
	"pertanyaan": "b.pertanyaan",
	"jawaban":    "a.jawaban",
}

// ------------------------
// GET ALL
// ------------------------
func (r *TemplateJawabanRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomaintemplatejawaban.SearchFilter,
	page, limit *int,
	deleted bool,
) ([]domaintemplatejawaban.TemplateJawabanDefault, int64, error) {

	var rows = make([]domaintemplatejawaban.TemplateJawabanDefault, 0)
	var total int64

	db := r.db.Debug().WithContext(ctx).
		Table("template_pilihanv2 a").
		Joins("LEFT JOIN template_pertanyaanv2 b ON b.id = a.id_template_pertanyaan").
		Select(`
		a.id as ID,
		a.uuid as UUID,
		a.id_template_pertanyaan as IdTemplatePertanyaan,
		b.uuid as UUIDTemplatePertanyaan,
		b.pertanyaan as NamaTemplatePertanyaan,
		a.jawaban as Jawaban,
		a.nilai as Nilai,
		a.isFreeText as IsFreeText,
		a.deleted_at as DeletedAt,
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
func (r *TemplateJawabanRepository) Create(ctx context.Context, templatejawaban *domaintemplatejawaban.TemplateJawaban) error {
	return r.db.WithContext(ctx).Create(templatejawaban).Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *TemplateJawabanRepository) Update(ctx context.Context, templatejawaban *domaintemplatejawaban.TemplateJawaban) error {
	return r.db.WithContext(ctx).Save(templatejawaban).Error
}

// ------------------------
// DELETE (by UUID)
// ------------------------
func (r *TemplateJawabanRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domaintemplatejawaban.TemplateJawaban{}).Error
}

func (r *TemplateJawabanRepository) SetupUuid(ctx context.Context) error {
	const chunkSize = 500

	var ids []uint
	if err := r.db.WithContext(ctx).
		Model(&domaintemplatejawaban.TemplateJawaban{}).
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
			"UPDATE template_pilihanv2 SET uuid = %s WHERE id IN (?)",
			caseSQL,
		)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}
