package infrastructure

import (
	domainkuesioner "UnpakSiamida/modules/kuesioner/domain"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type KuesionerJawabanRepository struct {
	db *gorm.DB
}

func NewKuesionerJawabanRepository(db *gorm.DB) domainkuesioner.IKuesionerJawabanRepository {
	return &KuesionerJawabanRepository{db: db}
}

// ===============================
// TX SUPPORT
// ===============================
func (r *KuesionerJawabanRepository) WithTx(tx any) domainkuesioner.IKuesionerJawabanRepository {
	return &KuesionerJawabanRepository{
		db: tx.(*gorm.DB),
	}
}

func (r *KuesionerJawabanRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	return r.db.WithContext(ctx).Begin(), nil
}

// ===============================
// GET ALL BY KUESIONER UUID
// ===============================
func (r *KuesionerJawabanRepository) GetAllByKuesioner(
	ctx context.Context,
	uuidkuesioner string,
) ([]domainkuesioner.KuesionerJawabanDefault, error) {

	results := make([]domainkuesioner.KuesionerJawabanDefault, 0)

	err := r.db.WithContext(ctx).
		Table("kuesioner_jawabanv2 kj").
		Select(`
			kj.id AS ID,
			kj.uuid AS UUID,

			kj.id_kuesioner AS IdKuesioner,
			k.uuid AS UuidKuesioner,

			kj.id_template_pertanyaan AS IdTemplatePertanyaan,
			tp.uuid AS UuidTemplatePertanyaan,

			kj.id_template_jawaban AS IdTemplateJawaban,
			tj.uuid AS UuidTemplateJawaban,

			kj.freeText AS FreeText
		`).
		Joins("JOIN kuesionerv2 k ON k.id = kj.id_kuesioner").
		Joins("LEFT JOIN template_pertanyaanv2 tp ON tp.id = kj.id_template_pertanyaan").
		Joins("LEFT JOIN template_pilihanv2 tj ON tj.id = kj.id_template_jawaban").
		Where("k.uuid = ?", uuidkuesioner).
		Order("kj.id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

// ===============================
// GET BY PERTANYAAN + USER
// ===============================
func (r *KuesionerJawabanRepository) GetByPertanyaanAndUser(
	ctx context.Context,
	pertanyaanID uint,
	sid string,
	resource string,
) ([]domainkuesioner.KuesionerJawaban, error) {

	results := make([]domainkuesioner.KuesionerJawaban, 0)

	err := r.db.WithContext(ctx).
		Where("id_template_pertanyaan = ?", pertanyaanID).
		Where("createdByRef = ?", sid).
		Where("createdBy = ?", resource).
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

// ===============================
// CREATE
// ===============================
func (r *KuesionerJawabanRepository) Create(
	ctx context.Context,
	data *domainkuesioner.KuesionerJawaban,
) error {

	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id_kuesioner"},
				{Name: "id_template_pertanyaan"},
				{Name: "id_template_jawaban"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"id_template_jawaban",
				"freeText",
				"updated_at",
			}),
		}).
		Create(data).Error
}

// ===============================
// DELETE
// ===============================
func (r *KuesionerJawabanRepository) Delete(
	ctx context.Context,
	id uint,
) error {

	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&domainkuesioner.KuesionerJawaban{}).Error
}

func (r *KuesionerJawabanRepository) GetTotalInputByKuesionerIDs(
	ctx context.Context,
	ids []uint,
) (map[string]uint, error) {

	type row struct {
		UuidKuesioner string
		Total         uint
	}

	var rows []row

	// Subquery: hitung 1 row per pertanyaan unik per kuesioner
	subQuery := r.db.
		Table("kuesioner_jawabanv2 kj").
		Select("k.uuid AS uuid").
		Joins("JOIN kuesionerv2 k ON kj.id_kuesioner = k.id").
		Joins("JOIN template_pilihanv2 tp on kj.id_template_jawaban = tp.id").
		Where("kj.id_kuesioner IN ?", ids).
		Where("tp.deleted_at IS NULL").
		Group("k.uuid, kj.id_template_pertanyaan")

	// Main query: hitung total pertanyaan per kuesioner
	err := r.db.WithContext(ctx).
		Debug().
		Table("(?) AS sub", subQuery).
		Select("uuid AS UuidKuesioner, COUNT(*) AS Total").
		Group("uuid").
		Having("COUNT(*) > 0").
		Find(&rows).Error

	if err != nil {
		return nil, err
	}

	// Mapping hasil ke map
	result := make(map[string]uint, len(rows))
	for _, r := range rows {
		result[r.UuidKuesioner] = r.Total
	}

	return result, nil
}
