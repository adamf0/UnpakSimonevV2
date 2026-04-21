package infrastructure

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"strings"

	"gorm.io/gorm"

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

func (r *AccountRepository) Auth(ctx context.Context, username string, password string) (*domain.Account, error) {

	// Rekursi chain function
	var chain func(func(ctx context.Context) (*domain.Account, error), ...func(ctx context.Context) (*domain.Account, error)) (*domain.Account, error)
	chain = func(current func(ctx context.Context) (*domain.Account, error), rest ...func(ctx context.Context) (*domain.Account, error)) (*domain.Account, error) {
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
		func(ctx context.Context) (*domain.Account, error) { return r.authDB(ctx, username, password) },
		func(ctx context.Context) (*domain.Account, error) { return r.authSimak(ctx, username, password) },
		func(ctx context.Context) (*domain.Account, error) { return r.authSimpeg(ctx, username, password) },
	)
}

func (r *AccountRepository) Get(ctx context.Context, id domain.AccountIdentifier) (*domain.Account, error) {
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

func (r *AccountRepository) authDB(ctx context.Context, username string, password string) (*domain.Account, error) {

	var user domain.Account

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

func (r *AccountRepository) authSimak(ctx context.Context, username string, password string) (*domain.Account, error) {

	var user domain.Account

	query := `
WITH dosen_cte AS (
    SELECT 
        d.nidn,
        CAST(d.nama_dosen AS CHAR(255)) AS Name,
        CAST(NULLIF(TRIM(d.email), '') AS CHAR(255)) AS Email,
        f.kode_fakultas AS RefFakultas,
        f.nama_fakultas AS Fakultas,
        p.kode_prodi AS RefProdi,
        p.nama_prodi AS Prodi
    FROM m_dosen d
    LEFT JOIN m_fakultas f ON f.kode_fakultas = d.kode_fak
    LEFT JOIN r_prodi p ON p.kode_prodi = d.kode_prodi
),
mahasiswa_cte AS (
    SELECT
        m.nim,
        CAST(m.nama_mahasiswa AS CHAR(255)) AS Name,
        CAST(NULLIF(TRIM(m.email), '') AS CHAR(255)) AS Email,
        f.kode_fakultas AS RefFakultas,
        f.nama_fakultas AS Fakultas,
        p.kode_prodi AS RefProdi,
        p.nama_prodi AS Prodi
    FROM m_mahasiswa m
    LEFT JOIN m_fakultas f ON f.kode_fakultas = m.kode_fak
    LEFT JOIN r_prodi p ON p.kode_prodi = m.kode_prodi
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

func (r *AccountRepository) authSimpeg(ctx context.Context, username string, password string) (*domain.Account, error) { //[review] hanya pegawai saja seharusnya

	var user domain.Account

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

func (r *AccountRepository) getDB(ctx context.Context, userid string) (*domain.Account, error) {

	var user domain.Account

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

func (r *AccountRepository) getSimakDosen(ctx context.Context, nidn string) (*domain.Account, error) {

	var user domain.Account

	query := `
WITH dosen_cte AS (
    SELECT 
        d.nidn,
        CAST(d.nama_dosen AS CHAR(255)) AS Name,
        CAST(NULLIF(TRIM(d.email), '') AS CHAR(255)) AS Email,
        f.kode_fakultas AS RefFakultas,
        f.nama_fakultas AS Fakultas,
        p.kode_prodi AS RefProdi,
        p.nama_prodi AS Prodi
    FROM m_dosen d
    LEFT JOIN m_fakultas f ON f.kode_fakultas = d.kode_fak
    LEFT JOIN r_prodi p ON p.kode_prodi = d.kode_prodi
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

func (r *AccountRepository) getSimakMahasiswa(ctx context.Context, nim string) (*domain.Account, error) {

	var user domain.Account

	query := `
WITH mahasiswa_cte AS (
    SELECT
        m.nim,
        CAST(m.nama_mahasiswa AS CHAR(255)) AS Name,
        CAST(NULLIF(TRIM(m.email), '') AS CHAR(255)) AS Email,
        f.kode_fakultas AS RefFakultas,
        f.nama_fakultas AS Fakultas,
        p.kode_prodi AS RefProdi,
        p.nama_prodi AS Prodi
    FROM m_mahasiswa m
    LEFT JOIN m_fakultas f ON f.kode_fakultas = m.kode_fak
    LEFT JOIN r_prodi p ON p.kode_prodi = m.kode_prodi
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

func (r *AccountRepository) getSimpeg(ctx context.Context, nip *string, nidn *string) (*domain.Account, error) {

	var user domain.Account

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

func md5String(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
