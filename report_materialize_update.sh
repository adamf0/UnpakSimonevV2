#!/bin/bash

# ==================================================
# DATABASE CONFIG
# ==================================================
DB_HOST="127.0.0.1"
DB_PORT=""
DB_NAME=""
DB_USER=""
DB_PASS=""

# ==================================================
# LOG CONFIG
# ==================================================
LOG_FILE="kuesioner_materialized.log"
ERROR_LOG="kuesioner_materialized_error.log"

echo "==================================================" >> "$LOG_FILE"
echo "[$(date '+%Y-%m-%d %H:%M:%S')] Refresh started" >> "$LOG_FILE"

MYSQL_OUTPUT=$(
mysql \
    -h"$DB_HOST" \
    -P"$DB_PORT" \
    -u"$DB_USER" \
    -p"$DB_PASS" \
    "$DB_NAME" 2>&1 <<'SQL'

SET FOREIGN_KEY_CHECKS=0;

-- hapus temp table jika masih ada dari proses sebelumnya
DROP TABLE IF EXISTS kuesioner_materialized_new;

-- clone struktur + index + partition
CREATE TABLE kuesioner_materialized_new
LIKE kuesioner_materialized;

-- isi data ke tabel baru
INSERT INTO kuesioner_materialized_new (
    ID,
    UUID,
    Tanggal,

    NIDN,
    NamaDosen,
    NIP,
    NamaTendik,

    NPM,
    NamaMahasiswa,

    KodeFakultas,
    Fakultas,

    KodeProdi,
    Prodi,
    Unit,

    Judul,
    Semester,

    Pertanyaan,
    Jawaban,
    FreeText,

    JenisPilihan,
    Kategori,
    FullPath,

    partition_key
)
SELECT
    k.id,
    k.uuid,
    k.tanggal,

    k.nidn,
    k.nama_dosen,
    k.nip,
    k.nama_tendik,

    k.npm,
    k.nama_mahasiswa,

    k.kode_fakultas,
    k.fakultas,

    k.kode_prodi,
    k.prodi,
    k.unit,

    b.judul,
    b.semester,

    tp.pertanyaan,
    tj.jawaban,
    kj.freeText,

    tp.jenis_pilihan,
    ka.nama_kategori,
    ka.full_text,

    CASE
        WHEN k.unit IS NOT NULL
             AND k.unit <> ''
        THEN 'UNIT'
        ELSE COALESCE(k.kode_fakultas, 'UNKNOWN')
    END AS partition_key

FROM kuesionerv2 k
JOIN bank_soalv2 b
    ON b.id = k.id_bank_soal
JOIN kuesioner_jawabanv2 kj
    ON kj.id_kuesioner = k.id
LEFT JOIN template_pertanyaanv2 tp
    ON tp.id = kj.id_template_pertanyaan
LEFT JOIN kategoriv2 ka
    ON ka.id = tp.id_kategori
LEFT JOIN template_pilihanv2 tj
    ON tj.id = kj.id_template_jawaban;

-- jika insert berhasil, baru swap

DROP TABLE IF EXISTS kuesioner_materialized_old;

RENAME TABLE
    kuesioner_materialized TO kuesioner_materialized_old,
    kuesioner_materialized_new TO kuesioner_materialized;

SET FOREIGN_KEY_CHECKS=1;

SQL
)

MYSQL_EXIT_CODE=$?

if [ $MYSQL_EXIT_CODE -eq 0 ]; then
    {
        echo "[$(date '+%Y-%m-%d %H:%M:%S')] Refresh success"
        echo
    } >> "$LOG_FILE"
else
    {
        echo "=================================================="
        echo "[$(date '+%Y-%m-%d %H:%M:%S')] Refresh FAILED"
        echo "Exit Code : $MYSQL_EXIT_CODE"
        echo
        echo "$MYSQL_OUTPUT"
        echo
    } >> "$ERROR_LOG"

    {
        echo "[$(date '+%Y-%m-%d %H:%M:%S')] Refresh FAILED"
        echo "See error log: $ERROR_LOG"
        echo
    } >> "$LOG_FILE"
fi