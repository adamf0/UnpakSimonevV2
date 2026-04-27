package infrastructure

import (
	commondomainfakultas "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	domainfakultas "UnpakSiamida/modules/fakultas/domain"
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type FakultasRepository struct {
	db *gorm.DB
}

func NewFakultasRepository(db *gorm.DB) domainfakultas.IFakultasRepository {
	return &FakultasRepository{db: db}
}

var allowedSearchColumns = map[string]string{
	"kode_fakultas": "kode_fakultas",
	"nama_fakultas": "nama_fakultas",
}

/*
	=========================
	  GET ALL PURE TABLE fakultas

=========================
*/
func (r *FakultasRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commondomainfakultas.SearchFilter,
	page, limit *int,
) ([]domainfakultas.FakultasDefault, int64, error) {

	var rows []domainfakultas.FakultasDefault
	var total int64

	db := r.db.WithContext(ctx).
		Table("fakultas").
		Select(`
			kode_fakultas as KodeFakultas,
			nama_fakultas as NamaFakultas
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
			kode_fakultas LIKE ?
			OR nama_fakultas LIKE ?
		`, like, like, like)
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
	db = db.Order("nama_fakultas ASC")

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
