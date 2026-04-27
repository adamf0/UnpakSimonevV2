package infrastructure

import (
	commondomainprodi "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	domainprodi "UnpakSiamida/modules/prodi/domain"
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type ProdiRepository struct {
	db *gorm.DB
}

func NewProdiRepository(db *gorm.DB) domainprodi.IProdiRepository {
	return &ProdiRepository{db: db}
}

var allowedSearchColumns = map[string]string{
	"kode_prodi":    "kode_prodi",
	"kode_fakultas": "kode_fak",
	"nama_prodi": `
		CONCAT(
			nama_prodi, 
			CASE kode_jenjang
				WHEN 'E' THEN ' (D3)'
				WHEN 'A' THEN ' (S3)'
				WHEN 'B' THEN ' (S2)'
				WHEN 'C' THEN ' (S1)'
				ELSE ''
			END
		)
	`,
}

/*
	=========================
	  GET ALL PURE TABLE prodi

=========================
*/
func (r *ProdiRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomainprodi.SearchFilter,
	page, limit *int,
) ([]domainprodi.ProdiDefault, int64, error) {

	var rows []domainprodi.ProdiDefault
	var total int64

	db := r.db.WithContext(ctx).
		Table("m_program_studi").
		Select(`
			kode_prodi as KodeProdi,
			kode_fak as KodeFakultas,
			CONCAT(
				nama_prodi, 
				CASE kode_jenjang
					WHEN 'E' THEN ' (D3)'
					WHEN 'A' THEN ' (S3)'
					WHEN 'B' THEN ' (S2)'
					WHEN 'C' THEN ' (S1)'
					ELSE ''
				END
			) as NamaProdi
		`)

	/* =========================
	   ADVANCED FILTER
	========================= */
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
			db = db.Where(
				fmt.Sprintf("%s LIKE ?", col),
				"%"+helper.EscapeLike(val)+"%",
			)

		case "in":
			rawVals := strings.Split(val, ",")
			vals := make([]string, 0)

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

	/* =========================
	   GLOBAL SEARCH
	========================= */
	if s := strings.TrimSpace(search); s != "" {
		like := "%" + helper.EscapeLike(s) + "%"

		db = db.Where(`
			kode_prodi LIKE ?
			OR kode_fak LIKE ?
			OR nama_prodi LIKE ?
		`, like, like, like, like, like, like)
	}

	/* =========================
	   COUNT
	========================= */
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	/* =========================
	   ORDERING
	========================= */
	db = db.Order("nama_prodi ASC")

	/* =========================
	   PAGINATION
	========================= */
	if page != nil && limit != nil && *limit > 0 {
		offset := (*page - 1) * (*limit)

		db = db.Offset(offset).Limit(*limit)
	}

	/* =========================
	   EXECUTE
	========================= */
	if err := db.Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}
