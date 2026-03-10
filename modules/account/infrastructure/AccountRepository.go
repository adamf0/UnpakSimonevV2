package infrastructure

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

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

	type result struct {
		priority int
		user     *domain.Account
	}

	g, ctx := errgroup.WithContext(ctx)
	results := make([]result, 0)
	mu := sync.Mutex{}

	dbs := []struct {
		priority int
		fn       func(ctx context.Context) (*domain.Account, error)
	}{
		{1, func(ctx context.Context) (*domain.Account, error) { return r.authDB(ctx, username, password) }},
		{2, func(ctx context.Context) (*domain.Account, error) { return r.authSimak(ctx, username, password) }},
		{3, func(ctx context.Context) (*domain.Account, error) { return r.authSimpeg(ctx, username, password) }},
	}

	for _, db := range dbs {
		// db := db
		g.Go(func() error {
			user, err := db.fn(ctx)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Record not found → jangan cancel, goroutine lain mungkin sukses
					return nil
				}
				// Error fatal → cancel semua goroutine
				return err
			}
			mu.Lock()
			if user.ID != "" {
				results = append(results, result{db.priority, user})
			}
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		// Salah satu goroutine error → cancel semua
		return nil, err
	}

	// Jika semua DB tidak menemukan user
	if len(results) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// Pilih user terbaik sesuai priority
	bestPriority := 999
	var bestUser *domain.Account
	for _, r := range results {
		if r.priority < bestPriority {
			bestPriority = r.priority
			bestUser = r.user
		}
	}

	return bestUser, nil
}

func (r *AccountRepository) Get(ctx context.Context, id domain.AccountIdentifier) (*domain.Account, error) {

	if id.NIDN != nil && *id.NIDN != "" {
		return r.getSimakDosen(ctx, *id.NIDN)
	}

	if id.NIM != nil && *id.NIM != "" {
		return r.getSimakMahasiswa(ctx, *id.NIM)
	}

	if id.NIP != nil && *id.NIP != "" {
		return r.getSimpeg(ctx, *id.NIP)
	}

	if id.UserID != nil && *id.UserID != "" {
		return r.getDB(ctx, *id.UserID)
	}

	return nil, errors.New("identifier not provided")
}

func (r *AccountRepository) authDB(ctx context.Context, username string, password string) (*domain.Account, error) {

	var user domain.Account

	err := r.db.Debug().WithContext(ctx).
		Table("users").
		Select(`"local" as Resource, users.*`).
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
	u.userid as ID,
	"simak" as Resource,
    u.username AS Username,
    u.password AS Password,
    LOWER(u.level) AS Level,
    COALESCE(d.Name, m.Name) AS Name,
    COALESCE(d.Email, m.Email) AS Email,
    COALESCE(d.RefFakultas, m.RefFakultas) AS RefFakultas,
    COALESCE(d.Fakultas, m.Fakultas) AS Fakultas,
    COALESCE(d.RefProdi, m.RefProdi) AS RefProdi,
    COALESCE(d.Prodi, m.Prodi) AS Prodi,
    NULL AS Unit
FROM user u
LEFT JOIN dosen_cte d ON d.nidn = u.userid
LEFT JOIN mahasiswa_cte m ON m.nim = u.userid
WHERE u.username = ? and u.password = ? and u.level in ("MAHASISWA","DOSEN") and u.aktif="Y"
LIMIT 1
`
	hash := sha1.Sum([]byte(md5String(password)))
	hashString := hex.EncodeToString(hash[:])

	err := r.dbSimak.Debug().WithContext(ctx).
		Raw(query, username, hashString).
		Scan(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AccountRepository) authSimpeg(ctx context.Context, username string, password string) (*domain.Account, error) {

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
	t.unit AS Unit
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
		Table("users").
		Where("uuid = ?", userid).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AccountRepository) getSimakDosen(ctx context.Context, nidn string) (*domain.Account, error) {

	var user domain.Account

	query := `
SELECT
	d.nidn AS Username,
	d.nama_dosen AS Name,
	d.email AS Email,
	f.kode_fakultas AS RefFakultas,
	f.nama_fakultas AS Fakultas,
	p.kode_prodi AS RefProdi,
	p.nama_prodi AS Prodi,
	NULL AS Unit
FROM m_dosen d
LEFT JOIN m_fakultas f ON f.kode_fakultas = d.kode_fak
LEFT JOIN r_prodi p ON p.kode_prodi = d.kode_prodi
WHERE d.nidn = ?
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
SELECT
	m.nim AS Username,
	m.nama_mahasiswa AS Name,
	m.email AS Email,
	f.kode_fakultas AS RefFakultas,
	f.nama_fakultas AS Fakultas,
	p.kode_prodi AS RefProdi,
	p.nama_prodi AS Prodi,
	NULL AS Unit
FROM m_mahasiswa m
LEFT JOIN m_fakultas f ON f.kode_fakultas = m.kode_fak
LEFT JOIN r_prodi p ON p.kode_prodi = m.kode_prodi
WHERE m.nim = ?
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

func (r *AccountRepository) getSimpeg(ctx context.Context, nip string) (*domain.Account, error) {

	var user domain.Account

	query := `
SELECT
	t.nip AS Username,
	t.nama AS Name,
	NULL AS Email,
	NULL AS RefFakultas,
	t.fakultas AS Fakultas,
	NULL AS RefProdi,
	NULL AS Prodi,
	t.unit AS Unit
FROM v_tendik t
WHERE t.nip = ?
LIMIT 1
`

	err := r.dbSimpeg.WithContext(ctx).
		Raw(query, nip).
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
