package infrastructure

import (
	commondomainkuesioner "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	domainkuesioner "UnpakSiamida/modules/kuesioner/domain"
	"context"
	"errors"
	"fmt"
	"strings"

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

	// Ambil hanya kolom yang benar-benar ada di struct KuesionerDefault
	var rowData domainkuesioner.KuesionerDefault

	err := r.db.WithContext(ctx).
		Table("bank_soalv2 a").
		Select(`
			id as ID,
			uuid as UUID,
			judul as Judul,
			content as Content,
			deskripsi as Deskripsi,
			semester as Semester,
			tanggal_mulai as TanggalMulai,
			tanggal_akhir as TanggalAkhir,
			created_by as CreatedBy,
			created_by_ref as CreatedByRef,
			deleted_at as DeletedAt,
			status as Status
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
	"judul":    "a.judul",
	"semester": "a.semeter",
	"status":   "a.status",
}

// ------------------------
// GET ALL
// ------------------------
func (r *KuesionerRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomainkuesioner.SearchFilter,
	page, limit *int,
	deleted bool,
) ([]domainkuesioner.KuesionerDefault, int64, error) {

	var rows = make([]domainkuesioner.KuesionerDefault, 0)
	var total int64

	db := r.db.Debug().WithContext(ctx).
		Table("bank_soalv2 a").
		Select(`
			id as ID,
			uuid as UUID,
			judul as Judul,
			content as Content,
			deskripsi as Deskripsi,
			semester as Semester,
			tanggal_mulai as TanggalMulai,
			tanggal_akhir as TanggalAkhir,
			createdBy as CreatedBy,
			createdByRef as CreatedByRef,
			deleted_at as DeletedAt,
			status as Status
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

// [pr] masih gagal dapat current uuid setelah upsert
// ------------------------
// CREATE
// ------------------------
func (r *KuesionerRepository) Create(ctx context.Context, kuesioner *domainkuesioner.Kuesioner) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id_bank_soal"},
				{Name: "nidn"},
				{Name: "nip"},
				{Name: "npm"},
			},
			DoUpdates: clause.Assignments(map[string]interface{}{
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
		Create(kuesioner).Error
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
