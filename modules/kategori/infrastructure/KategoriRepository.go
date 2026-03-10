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

	err := r.db.WithContext(ctx).
		Table("kategori a").
		Select(`
			id as ID,
			uuid as UUID,
			nama_kategori as NamaKategori,
			full_text as FullTexts,
			sub_kategori as SubKategori,
			createdBy as CreatedBy,
			createdByRef as CreatedByRef,
			deleted_at as DeletedAt
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
	"kategori": "a.nama_kategori",
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

	db := r.db.Debug().WithContext(ctx).
		Table("kategori a").
		Select(`
			id as ID,
			uuid as UUID,
			nama_kategori as NamaKategori,
			full_text as FullTexts,
			sub_kategori as SubKategori,
			createdBy as CreatedBy,
			createdByRef as CreatedByRef,
			deleted_at as DeletedAt
	`)

	if deleted {
		db = db.Where(clause.Expr{
			SQL: "deleted_at IS NOT NULL",
		})
	} else {
		db = db.Where(clause.Expr{
			SQL: "deleted_at IS NULL",
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

	var rows []domainkategori.Kategori

	err := r.db.WithContext(ctx).
		Where("sub_kategori = ?", parentID).
		Find(&rows).Error

	return rows, err
}

func (r *KategoriRepository) RebuildFullText(ctx context.Context) error {

	var list []domainkategori.Kategori

	if err := r.db.WithContext(ctx).
		Find(&list).Error; err != nil {
		return err
	}

	childrenMap := map[uint][]domainkategori.Kategori{}
	var roots []domainkategori.Kategori

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
