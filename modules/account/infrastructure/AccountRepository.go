package infrastructure

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	commonDomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	domain "UnpakSiamida/modules/account/domain"
)

type AccountRepository struct {
	db       *gorm.DB
	dbSimak  *gorm.DB
	dbSimpeg *gorm.DB
}

func NewAccountRepository(db *gorm.DB, dbsimak *gorm.DB, dbsimpeg *gorm.DB) domain.IAccountRepository {
	return &AccountRepository{
		db:       db,
		dbSimak:  dbsimak,
		dbSimpeg: dbsimpeg,
	}
}

func (r *AccountRepository) Auth(ctx context.Context, username string, password string) (*domain.AccountDefault, error) {

	// Rekursi chain function
	var chain func(func(ctx context.Context) (*domain.AccountDefault, error), ...func(ctx context.Context) (*domain.AccountDefault, error)) (*domain.AccountDefault, error)
	chain = func(current func(ctx context.Context) (*domain.AccountDefault, error), rest ...func(ctx context.Context) (*domain.AccountDefault, error)) (*domain.AccountDefault, error) {
		user, err := current(ctx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Not found → lanjut ke function berikutnya
				if len(rest) == 0 {
					return nil, gorm.ErrRecordNotFound
				}
				return chain(rest[0], rest[1:]...)
			}
			// Fatal error → stop
			return nil, err
		}

		if user != nil && user.ID != "" {
			return user, nil
		}

		if len(rest) == 0 {
			return nil, gorm.ErrRecordNotFound
		}
		return chain(rest[0], rest[1:]...)
	}

	// Panggil chain sesuai urutan prioritas: authDB -> authSimak -> authSimpeg
	return chain(
		func(ctx context.Context) (*domain.AccountDefault, error) { return r.authDB(ctx, username, password) },
		func(ctx context.Context) (*domain.AccountDefault, error) { return r.authSimak(ctx, username, password) },
		func(ctx context.Context) (*domain.AccountDefault, error) {
			return r.authSimpeg(ctx, username, password)
		},
	)
}

func (r *AccountRepository) Get(ctx context.Context, id domain.AccountIdentifier) (*domain.AccountDefault, error) {
	if strings.TrimSpace(helper.StringValue(id.NIDN)) != "" {
		return r.getSimakDosen(ctx, helper.StringValue(id.NIDN))
	}

	if strings.TrimSpace(helper.StringValue(id.NIM)) != "" {
		return r.getSimakMahasiswa(ctx, helper.StringValue(id.NIM))
	}

	if strings.TrimSpace(helper.StringValue(id.NIP)) != "" {
		return r.getSimpeg(ctx, id.NIP, id.NIDN)
	}

	if strings.TrimSpace(helper.StringValue(id.UserID)) != "" {
		return r.getDB(ctx, helper.StringValue(id.UserID))
	}

	return nil, errors.New("identifier not provided")
}

func (r *AccountRepository) authDB(ctx context.Context, username string, password string) (*domain.AccountDefault, error) {

	var user domain.AccountDefault

	err := r.db.Debug().WithContext(ctx).
		Table("users").
		Select(`"local" as Resource, null as CodeCtx, users.*`).
		Where("username = ? AND password_plain = ?", username, password).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AccountRepository) authSimak(ctx context.Context, username string, password string) (*domain.AccountDefault, error) {

	var user domain.AccountDefault

	query := `
WITH dosen_cte AS (
    SELECT 
        d.nidn,
        CAST(d.nama_dosen AS CHAR(255)) AS Name,
        CAST(NULLIF(TRIM(d.email), '') AS CHAR(255)) AS Email,
        f.kode_fakultas AS RefFakultas,
        f.nama_fakultas AS Fakultas,
        p.kode_prodi AS RefProdi,
        CONCAT(
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
		) AS Prodi 
    FROM m_dosen d
    LEFT JOIN m_fakultas f ON f.kode_fakultas = d.kode_fak
    LEFT JOIN m_program_studi p ON p.kode_prodi = d.kode_prodi
),
mahasiswa_cte AS (
    SELECT
        m.nim,
        CAST(m.nama_mahasiswa AS CHAR(255)) AS Name,
        CAST(NULLIF(TRIM(m.email), '') AS CHAR(255)) AS Email,
        f.kode_fakultas AS RefFakultas,
        f.nama_fakultas AS Fakultas,
        p.kode_prodi AS RefProdi,
        CONCAT(
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
		) AS Prodi 
    FROM m_mahasiswa m
    LEFT JOIN m_fakultas f ON f.kode_fakultas = m.kode_fak
    LEFT JOIN m_program_studi p ON p.kode_prodi = m.kode_prodi
)
SELECT
	u.userid AS ID,
	"simak" AS Resource,
    u.username AS Username,
    u.password AS Password,
    LOWER(u.level) AS Level,
    COALESCE(d.Name, m.Name) AS Name,
    COALESCE(d.Email, m.Email) AS Email,
    COALESCE(d.RefFakultas, m.RefFakultas) AS RefFakultas,
    COALESCE(d.Fakultas, m.Fakultas) AS Fakultas,
    COALESCE(d.RefProdi, m.RefProdi) AS RefProdi,
    COALESCE(d.Prodi, m.Prodi) AS Prodi,
    NULL AS Unit,
	CASE 
		WHEN LENGTH(TRIM(d.Name)) > 0 THEN ?
		WHEN LENGTH(TRIM(m.Name)) > 0 THEN ?
		ELSE ''
	END AS CodeCtx
FROM user u
LEFT JOIN dosen_cte d ON d.nidn = u.userid
LEFT JOIN mahasiswa_cte m ON m.nim = u.userid
WHERE 
	u.username = ?
	AND u.password = ?
	AND u.level IN ("MAHASISWA","DOSEN")
	AND u.aktif = "Y"
LIMIT 1;
`
	hash := sha1.Sum([]byte(md5String(password)))
	hashString := hex.EncodeToString(hash[:])

	err := r.dbSimak.Debug().WithContext(ctx).
		Raw(query, domain.CtxDosen, domain.CtxMahasiswa, username, hashString).
		Scan(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AccountRepository) authSimpeg(ctx context.Context, username string, password string) (*domain.AccountDefault, error) { //[review] hanya pegawai saja seharusnya

	var user domain.AccountDefault

	query := `
SELECT
	u.id as ID,
	"simpeg" as Resource,
	u.username AS Username,
	u.password AS Password,
	LOWER(u.level) AS Level,
	t.nama AS Name,
	NULL AS Email,
	NULL AS RefFakultas,
	t.fakultas AS Fakultas,
	NULL AS RefProdi,
	NULL AS Prodi,
	t.unit AS Unit,
	null as CodeCtx
FROM pengguna u
LEFT JOIN v_tendik t ON t.nip = u.username
WHERE u.username = ? and u.password = ? and u.level in ("PEGAWAI","DOSEN") and u.status="AKTIF"
LIMIT 1
`
	hash := sha1.Sum([]byte(password))
	hashString := hex.EncodeToString(hash[:])

	err := r.dbSimpeg.Debug().WithContext(ctx).
		Raw(query, username, hashString).
		Scan(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AccountRepository) getDB(ctx context.Context, userid string) (*domain.AccountDefault, error) {

	var user domain.AccountDefault

	err := r.db.WithContext(ctx).
		Table("users u").
		Select(`
			u.*,
			u.fakultas as RefFakultas,
			f.nama_fakultas as Fakultas,
			u.prodi as RefProdi,
			p.nama_prodi as Prodi,
			null as Unit,
			'local' as Resource,
			null as CodeCtx
		`).
		Joins("LEFT JOIN m_fakultas f ON f.kode_fakultas = u.fakultas").
		Joins("LEFT JOIN m_program_studi p ON p.kode_prodi = u.prodi").
		Where("u.id = ?", userid).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AccountRepository) getSimakDosen(ctx context.Context, nidn string) (*domain.AccountDefault, error) {

	var user domain.AccountDefault

	query := `
WITH dosen_cte AS (
    SELECT 
        d.nidn,
        CAST(d.nama_dosen AS CHAR(255)) AS Name,
        CAST(NULLIF(TRIM(d.email), '') AS CHAR(255)) AS Email,
        f.kode_fakultas AS RefFakultas,
        f.nama_fakultas AS Fakultas,
        p.kode_prodi AS RefProdi,
        CONCAT(
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
		) AS Prodi 
    FROM m_dosen d
    LEFT JOIN m_fakultas f ON f.kode_fakultas = d.kode_fak
    LEFT JOIN m_program_studi p ON p.kode_prodi = d.kode_prodi
) 
SELECT
	u.userid as ID,
	"simak" as Resource,
    u.username AS Username,
    u.password AS Password,
    LOWER(u.level) AS Level,
    d.Name AS Name,
    d.Email AS Email,
    d.RefFakultas AS RefFakultas,
    d.Fakultas AS Fakultas,
    d.RefProdi AS RefProdi,
    d.Prodi AS Prodi,
    NULL AS Unit,
	'` + domain.CtxDosen + `' AS CodeCtx
FROM user u
LEFT JOIN dosen_cte d ON d.nidn = u.userid
WHERE d.nidn = ? and u.level = "DOSEN"
LIMIT 1
`

	err := r.dbSimak.WithContext(ctx).
		Raw(query, nidn).
		Scan(&user).Error

	if err != nil {
		return nil, err
	}

	// user.Level = "dosen"

	return &user, nil
}

func (r *AccountRepository) getSimakMahasiswa(ctx context.Context, nim string) (*domain.AccountDefault, error) {

	var user domain.AccountDefault

	query := `
WITH mahasiswa_cte AS (
    SELECT
        m.nim,
        CAST(m.nama_mahasiswa AS CHAR(255)) AS Name,
        CAST(NULLIF(TRIM(m.email), '') AS CHAR(255)) AS Email,
        f.kode_fakultas AS RefFakultas,
        f.nama_fakultas AS Fakultas,
        p.kode_prodi AS RefProdi,
        CONCAT(
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
		) AS Prodi 
    FROM m_mahasiswa m
    LEFT JOIN m_fakultas f ON f.kode_fakultas = m.kode_fak
    LEFT JOIN m_program_studi p ON p.kode_prodi = m.kode_prodi
)
SELECT
	u.userid as ID,
	"simak" as Resource,
    u.username AS Username,
    u.password AS Password,
    LOWER(u.level) AS Level,
    m.Name AS Name,
    m.Email AS Email,
    m.RefFakultas AS RefFakultas,
    m.Fakultas AS Fakultas,
    m.RefProdi AS RefProdi,
    m.Prodi AS Prodi,
    NULL AS Unit,
	'` + domain.CtxMahasiswa + `' AS CodeCtx
FROM user u
LEFT JOIN mahasiswa_cte m ON m.nim = u.userid
WHERE m.nim = ? and u.level = "MAHASISWA"
LIMIT 1
`

	err := r.dbSimak.WithContext(ctx).
		Raw(query, nim).
		Scan(&user).Error

	if err != nil {
		return nil, err
	}

	// user.Level = "mahasiswa"

	return &user, nil
}

func (r *AccountRepository) getSimpeg(ctx context.Context, nip *string, nidn *string) (*domain.AccountDefault, error) {

	var user domain.AccountDefault

	query := `
SELECT
	u.id as ID,
	"simpeg" as Resource,
	u.username AS Username,
	u.password AS Password,
	LOWER(u.level) AS Level,
	t.nama AS Name,
	NULL AS Email,
	NULL AS RefFakultas,
	t.fakultas AS Fakultas,
	NULL AS RefProdi,
	NULL AS Prodi,
	t.unit AS Unit,
	null as CodeCtx
FROM pengguna u
LEFT JOIN v_tendik t ON t.nip = u.username
WHERE u.id = ?
LIMIT 1
`

	err := r.dbSimpeg.WithContext(ctx).
		Raw(query, helper.StringValue(nip)).
		Scan(&user).Error

	if err != nil {
		return nil, err
	}

	// user.Level = "tendik"

	return &user, nil
}

// ------------------------
// GET BY UUID
// ------------------------
func (r *AccountRepository) GetByUuid(
	ctx context.Context,
	uid uuid.UUID,
) (*domain.Account, error) {

	var account domain.Account

	err := r.db.WithContext(ctx).
		Debug().
		Table("users u").
		Select(`
			u.id as ID,
			u.uuid as UUID,
			u.username as Username,
			u.password as Password,
			u.level as Level,
			u.name as Name,
			u.email as Email,
			u.fakultas as RefFakultas,
			f.nama_fakultas as Fakultas,
			u.prodi as RefProdi,
			p.nama_prodi as Prodi,
			u.deleted_at as DeletedAt,
			u.created_at as CreatedAt,
			u.updated_at as UpdatedAt
		`).
		Joins("LEFT JOIN m_fakultas f ON f.kode_fakultas = u.fakultas").
		Joins("LEFT JOIN m_program_studi p ON p.kode_prodi = u.prodi").
		Where("u.uuid = ?", uid).
		Take(&account).Error

	if err != nil {
		return nil, err
	}

	return &account, nil
}

var allowedSearchColumns = map[string]string{
	"username":      "u.username",
	"name":          "u.name",
	"email":         "u.email",
	"level":         "u.level",
	"fakultas":      "u.fakultas",
	"prodi":         "u.prodi",
	"nama_fakultas": "f.nama_fakultas",
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
}

// ------------------------
// GET ALL
// ------------------------
func (r *AccountRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commonDomain.SearchFilter,
	page, limit *int,
	deleted bool,
) ([]domain.Account, int64, error) {

	var rows []domain.Account
	var total int64

	db := r.db.WithContext(ctx).
		Table("users u").
		Select(`
			u.id as ID,
			u.uuid as UUID,
			u.username as Username,
			u.level as Level,
			u.name as Name,
			u.email as Email,
			u.fakultas as RefFakultas,
			f.nama_fakultas as Fakultas,
			u.prodi as RefProdi,
			CONCAT(
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
			) AS Prodi, 
			u.deleted_at as DeletedAt,
			u.created_at as CreatedAt,
			u.updated_at as UpdatedAt
		`).
		Joins("LEFT JOIN m_fakultas f ON f.kode_fakultas = u.fakultas").
		Joins("LEFT JOIN m_program_studi p ON p.kode_prodi = u.prodi")

	if deleted {
		db = db.Where("u.deleted_at IS NOT NULL")
	} else {
		db = db.Where("u.deleted_at IS NULL")
	}

	// ------------------------
	// ADVANCED FILTER
	// ------------------------
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
			vals := make([]interface{}, 0)

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

	// ------------------------
	// GLOBAL SEARCH
	// ------------------------
	if s := strings.TrimSpace(search); s != "" {
		like := "%" + helper.EscapeLike(s) + "%"

		var conditions []string
		var args []interface{}

		for _, col := range allowedSearchColumns {
			conditions = append(
				conditions,
				fmt.Sprintf("%s LIKE ?", col),
			)
			args = append(args, like)
		}

		db = db.Where(
			"("+strings.Join(conditions, " OR ")+")",
			args...,
		)
	}

	// ------------------------
	// COUNT
	// ------------------------
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ------------------------
	// PAGINATION
	// ------------------------
	db = db.Order("u.id DESC")

	if page != nil && limit != nil && *limit > 0 {
		offset := (*page - 1) * (*limit)
		db = db.Offset(offset).Limit(*limit)
	}

	// ------------------------
	// EXECUTE
	// ------------------------
	if err := db.Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

// ------------------------
// CREATE
// ------------------------
func (r *AccountRepository) Create(
	ctx context.Context,
	account *domain.Account,
) error {
	return r.db.WithContext(ctx).
		Create(account).Error
}

// ------------------------
// UPDATE
// ------------------------
func (r *AccountRepository) Update(
	ctx context.Context,
	account *domain.Account,
) error {
	return r.db.WithContext(ctx).
		Save(account).Error
}

// ------------------------
// DELETE
// ------------------------
func (r *AccountRepository) Delete(
	ctx context.Context,
	uid uuid.UUID,
) error {
	return r.db.WithContext(ctx).
		Where("uuid = ?", uid).
		Delete(&domain.Account{}).Error
}

// ------------------------
// SETUP UUID
// ------------------------
func (r *AccountRepository) SetupUuid(
	ctx context.Context,
) error {

	const chunkSize = 500

	var ids []string

	if err := r.db.WithContext(ctx).
		Model(&domain.Account{}).
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
		args := make([]any, 0)

		for _, id := range chunk {
			caseSQL += "WHEN ? THEN ? "
			args = append(args, id, uuid.NewString())
		}

		caseSQL += "END"

		args = append(args, chunk)

		query := fmt.Sprintf(`
			UPDATE users
			SET uuid = %s
			WHERE id IN (?)
		`, caseSQL)

		if err := r.db.WithContext(ctx).
			Exec(query, args...).Error; err != nil {
			return err
		}
	}

	return nil
}

func md5String(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
