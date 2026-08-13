#!/bin/bash

# ==============================================================================
# SIMONEV - BATCH REPORT SUMMARY CALCULATOR SCRIPT (CORRECTED FAKULTAS, PRODI & UNIT)
# ==============================================================================
# Script ini 100% pure .sh / python3 runner (tanpa kode Golang / .go file).
# Mengkalkulasi data dari kuesioner_materialized ke 5 tabel summary secara presisi:
# 1. report_summary_overview    (Overview per Overall, Fakultas, Prodi, & Unit)
# 2. report_distribusi_fakultas (Distribusi Fakultas dengan KodeFakultas & Unit)
# 3. report_top_questions       (Top 10 Questions per Filter)
# 4. report_kategori_summary     (Rating & JSON Questions per Filter)
# 5. report_year                 (Grafik 4 Year Chart)
# ==============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="${SCRIPT_DIR}/report_summary_calculator.log"
ERROR_LOG="${SCRIPT_DIR}/report_summary_calculator_error.log"

echo "==================================================" >> "$LOG_FILE"
echo "[$(date '+%Y-%m-%d %H:%M:%S')] Start pure bash report summary calculation" >> "$LOG_FILE"

export DB_HOST="${DB_HOST:-127.0.0.1}"
export DB_PORT="${DB_PORT:-3306}"
export DB_NAME="${DB_NAME:-unpak_simonevv2}"
export DB_USER="${DB_USER:-root}"
export DB_PASS="${DB_PASS:-}"
export MYSQL_BIN="$(which mysql 2>/dev/null || echo "/Applications/XAMPP/xamppfiles/bin/mysql")"

python3 - << 'PYTHON_SCRIPT'
import os
import sys
import subprocess
import json

host = os.environ.get("DB_HOST", "127.0.0.1")
port = os.environ.get("DB_PORT", "3306")
name = os.environ.get("DB_NAME", "unpak_simonevv2")
user = os.environ.get("DB_USER", "root")
password = os.environ.get("DB_PASS", "")
mysql_bin = os.environ.get("MYSQL_BIN", "mysql")

def run_sql(sql):
    cmd = [mysql_bin]
    if host:
        cmd.append(f"-h{host}")
    if port:
        cmd.append(f"-P{port}")
    if user:
        cmd.append(f"-u{user}")
    if password:
        cmd.append(f"-p{password}")
    cmd.extend([name, "-N", "-B"])
    
    p = subprocess.run(cmd, input=sql, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    if p.returncode != 0:
        raise Exception(f"SQL Error: {p.stderr.strip()}")
    return p.stdout.strip()

def escape_sql(val):
    if val is None:
        return "NULL"
    s = str(val).replace("\\", "\\\\").replace("'", "\\'")
    return f"'{s}'"

print("=== STARTING PRECISE BATCH SUMMARY CALCULATOR ===")

# Create tables if not exist
ddls = [
    "DROP TABLE IF EXISTS report_summary_overview",
    "DROP TABLE IF EXISTS report_distribusi_fakultas",
    "DROP TABLE IF EXISTS report_top_questions",
    "DROP TABLE IF EXISTS report_kategori_summary",
    "DROP TABLE IF EXISTS report_year",
    """CREATE TABLE IF NOT EXISTS report_summary_overview (
        id INT AUTO_INCREMENT PRIMARY KEY,
        judul VARCHAR(255) NOT NULL,
        semester VARCHAR(50) NOT NULL,
        kode_fakultas VARCHAR(50) NOT NULL DEFAULT '',
        kode_prodi VARCHAR(50) NOT NULL DEFAULT '',
        unit VARCHAR(255) NOT NULL DEFAULT '',
        total_responden INT DEFAULT 0,
        total_jawaban INT DEFAULT 0,
        rata_rata_rating DECIMAL(4,2) DEFAULT 0.00,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        KEY idx_judul_sem_org (judul(150), semester, kode_fakultas, kode_prodi, unit(100))
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4""",
    """CREATE TABLE IF NOT EXISTS report_distribusi_fakultas (
        id INT AUTO_INCREMENT PRIMARY KEY,
        judul VARCHAR(255) NOT NULL,
        semester VARCHAR(50) NOT NULL,
        kode_fakultas VARCHAR(50) NOT NULL DEFAULT '',
        kode_prodi VARCHAR(50) NOT NULL DEFAULT '',
        unit VARCHAR(255) NOT NULL DEFAULT '',
        nama_fakultas VARCHAR(255),
        total_responden INT DEFAULT 0,
        persentase DECIMAL(10,2) DEFAULT 0.00,
        prodi_distribution JSON,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        KEY idx_judul_sem_fak_prodi (judul(150), semester, kode_fakultas, kode_prodi, unit(100))
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4""",
    """CREATE TABLE IF NOT EXISTS report_top_questions (
        id INT AUTO_INCREMENT PRIMARY KEY,
        judul VARCHAR(255) NOT NULL,
        semester VARCHAR(50) NOT NULL,
        kode_fakultas VARCHAR(50) NOT NULL DEFAULT '',
        kode_prodi VARCHAR(50) NOT NULL DEFAULT '',
        unit VARCHAR(255) NOT NULL DEFAULT '',
        pertanyaan TEXT NOT NULL,
        nama_kategori VARCHAR(255),
        total_responden INT DEFAULT 0,
        rata_rata_skor DECIMAL(4,2) DEFAULT 0.00,
        peringkat INT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        KEY idx_judul_sem_rank (judul(150), semester, kode_fakultas, kode_prodi, unit(100), peringkat)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4""",
    """CREATE TABLE IF NOT EXISTS report_kategori_summary (
        id INT AUTO_INCREMENT PRIMARY KEY,
        judul VARCHAR(255) NOT NULL,
        semester VARCHAR(50) NOT NULL,
        kode_fakultas VARCHAR(50) NOT NULL DEFAULT '',
        kode_prodi VARCHAR(50) NOT NULL DEFAULT '',
        unit VARCHAR(255) NOT NULL DEFAULT '',
        nama_kategori VARCHAR(255) NOT NULL,
        full_text TEXT,
        total_pertanyaan INT DEFAULT 0,
        total_responden INT DEFAULT 0,
        rata_rata_skor DECIMAL(4,2) DEFAULT 0.00,
        chart_distribution JSON,
        questions_json JSON,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        KEY idx_judul_sem_kat (judul(150), semester, kode_fakultas, kode_prodi, unit(100), nama_kategori(100))
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4""",
    """CREATE TABLE IF NOT EXISTS report_year (
        id INT AUTO_INCREMENT PRIMARY KEY,
        tahun VARCHAR(50) NOT NULL,
        total_kuesioner INT DEFAULT 0,
        total_mahasiswa INT DEFAULT 0,
        total_dosen INT DEFAULT 0,
        total_tendik INT DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        UNIQUE KEY uq_tahun (tahun)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"""
]

for ddl in ddls:
    run_sql(ddl)

# 1. Populate report_summary_overview
# a) Overall
run_sql("""
INSERT INTO report_summary_overview (judul, semester, kode_fakultas, kode_prodi, unit, total_responden, total_jawaban, rata_rata_rating)
SELECT 
    k.Judul AS judul,
    CAST(k.Semester AS CHAR) AS semester,
    '' AS kode_fakultas,
    '' AS kode_prodi,
    '' AS unit,
    COUNT(DISTINCT COALESCE(NULLIF(k.NPM,''), NULLIF(k.NIDN,''), NULLIF(k.NIP,''), k.NamaMahasiswa, k.NamaDosen, k.NamaTendik)) AS total_responden,
    COUNT(*) AS total_jawaban,
    COALESCE(ROUND(AVG(CASE WHEN k.JenisPilihan = 'rating' AND k.Jawaban REGEXP '^[0-9]+$' THEN CAST(k.Jawaban AS DECIMAL(4,2)) ELSE NULL END), 2), 0.00) AS rata_rata_rating
FROM kuesioner_materialized k
WHERE k.Judul IS NOT NULL AND k.Judul != ''
GROUP BY k.Judul, k.Semester;
""")

# b) Per Fakultas
run_sql("""
INSERT INTO report_summary_overview (judul, semester, kode_fakultas, kode_prodi, unit, total_responden, total_jawaban, rata_rata_rating)
SELECT 
    k.Judul AS judul,
    CAST(k.Semester AS CHAR) AS semester,
    k.KodeFakultas AS kode_fakultas,
    '' AS kode_prodi,
    '' AS unit,
    COUNT(DISTINCT COALESCE(NULLIF(k.NPM,''), NULLIF(k.NIDN,''), NULLIF(k.NIP,''), k.NamaMahasiswa, k.NamaDosen, k.NamaTendik)) AS total_responden,
    COUNT(*) AS total_jawaban,
    COALESCE(ROUND(AVG(CASE WHEN k.JenisPilihan = 'rating' AND k.Jawaban REGEXP '^[0-9]+$' THEN CAST(k.Jawaban AS DECIMAL(4,2)) ELSE NULL END), 2), 0.00) AS rata_rata_rating
FROM kuesioner_materialized k
WHERE k.Judul IS NOT NULL AND k.Judul != '' AND k.KodeFakultas IS NOT NULL AND k.KodeFakultas != ''
GROUP BY k.Judul, k.Semester, k.KodeFakultas;
""")

# c) Per Prodi
run_sql("""
INSERT INTO report_summary_overview (judul, semester, kode_fakultas, kode_prodi, unit, total_responden, total_jawaban, rata_rata_rating)
SELECT 
    k.Judul AS judul,
    CAST(k.Semester AS CHAR) AS semester,
    k.KodeFakultas AS kode_fakultas,
    k.KodeProdi AS kode_prodi,
    '' AS unit,
    COUNT(DISTINCT COALESCE(NULLIF(k.NPM,''), NULLIF(k.NIDN,''), NULLIF(k.NIP,''), k.NamaMahasiswa, k.NamaDosen, k.NamaTendik)) AS total_responden,
    COUNT(*) AS total_jawaban,
    COALESCE(ROUND(AVG(CASE WHEN k.JenisPilihan = 'rating' AND k.Jawaban REGEXP '^[0-9]+$' THEN CAST(k.Jawaban AS DECIMAL(4,2)) ELSE NULL END), 2), 0.00) AS rata_rata_rating
FROM kuesioner_materialized k
WHERE k.Judul IS NOT NULL AND k.Judul != '' AND k.KodeProdi IS NOT NULL AND k.KodeProdi != ''
GROUP BY k.Judul, k.Semester, k.KodeFakultas, k.KodeProdi;
""")

# d) Per Unit (Unit respondents have empty KodeFakultas & KodeProdi)
run_sql("""
INSERT INTO report_summary_overview (judul, semester, kode_fakultas, kode_prodi, unit, total_responden, total_jawaban, rata_rata_rating)
SELECT 
    k.Judul AS judul,
    CAST(k.Semester AS CHAR) AS semester,
    '' AS kode_fakultas,
    '' AS kode_prodi,
    COALESCE(NULLIF(k.Unit, ''), 'Umum') AS unit,
    COUNT(DISTINCT COALESCE(NULLIF(k.NPM,''), NULLIF(k.NIDN,''), NULLIF(k.NIP,''), k.NamaMahasiswa, k.NamaDosen, k.NamaTendik)) AS total_responden,
    COUNT(*) AS total_jawaban,
    COALESCE(ROUND(AVG(CASE WHEN k.JenisPilihan = 'rating' AND k.Jawaban REGEXP '^[0-9]+$' THEN CAST(k.Jawaban AS DECIMAL(4,2)) ELSE NULL END), 2), 0.00) AS rata_rata_rating
FROM kuesioner_materialized k
WHERE k.Judul IS NOT NULL AND k.Judul != '' AND (k.KodeFakultas IS NULL OR k.KodeFakultas = '')
GROUP BY k.Judul, k.Semester, COALESCE(NULLIF(k.Unit, ''), 'Umum');
""")
print("1. report_summary_overview populated!")

# 2. Populate report_distribusi_fakultas
# a) Level 1 Academic (Per Fakultas)
run_sql("""
INSERT INTO report_distribusi_fakultas (judul, semester, kode_fakultas, kode_prodi, unit, nama_fakultas, total_responden, persentase, prodi_distribution)
SELECT 
    sub.Judul AS judul,
    CAST(sub.Semester AS CHAR) AS semester,
    sub.KodeFakultas AS kode_fakultas,
    '' AS kode_prodi,
    '' AS unit,
    COALESCE(sub.Fakultas, '') AS nama_fakultas,
    SUM(sub.prodi_resp) AS total_responden,
    0.00 AS persentase,
    CONCAT('[', COALESCE(GROUP_CONCAT(CONCAT('{"title":"', REPLACE(sub.n_prodi, '"', '\\\\\"'), '","total":', sub.prodi_resp, '}') ORDER BY sub.min_id ASC SEPARATOR ','), ''), ']') AS prodi_distribution
FROM (
    SELECT 
        k.Judul,
        k.Semester,
        k.KodeFakultas,
        k.Fakultas,
        COALESCE(NULLIF(k.Prodi, ''), 'Umum') AS n_prodi,
        COUNT(DISTINCT COALESCE(NULLIF(k.NPM,''), NULLIF(k.NIDN,''), NULLIF(k.NIP,''), k.NamaMahasiswa, k.NamaDosen, k.NamaTendik)) AS prodi_resp,
        MIN(k.id) AS min_id
    FROM kuesioner_materialized k
    WHERE k.Judul IS NOT NULL AND k.Judul != '' AND k.KodeFakultas IS NOT NULL AND k.KodeFakultas != ''
    GROUP BY k.Judul, k.Semester, k.KodeFakultas, k.Fakultas, COALESCE(NULLIF(k.Prodi, ''), 'Umum')
) sub
GROUP BY sub.Judul, sub.Semester, sub.KodeFakultas, sub.Fakultas;
""")

# b) Level 1 Unit (Per Unit for Tendik/Non-academic)
run_sql("""
INSERT INTO report_distribusi_fakultas (judul, semester, kode_fakultas, kode_prodi, unit, nama_fakultas, total_responden, persentase, prodi_distribution)
SELECT 
    k.Judul AS judul,
    CAST(k.Semester AS CHAR) AS semester,
    '' AS kode_fakultas,
    '' AS kode_prodi,
    COALESCE(NULLIF(k.Unit, ''), 'Umum') AS unit,
    COALESCE(k.Fakultas, '') AS nama_fakultas,
    COUNT(DISTINCT COALESCE(NULLIF(k.NPM,''), NULLIF(k.NIDN,''), NULLIF(k.NIP,''), k.NamaMahasiswa, k.NamaDosen, k.NamaTendik)) AS total_responden,
    0.00 AS persentase,
    CONCAT('[{"title":"', REPLACE(COALESCE(NULLIF(k.Prodi, ''), 'Umum'), '"', '\\\\\"'), '","total":', COUNT(DISTINCT COALESCE(NULLIF(k.NPM,''), NULLIF(k.NIDN,''), NULLIF(k.NIP,''), k.NamaMahasiswa, k.NamaDosen, k.NamaTendik)), '}]') AS prodi_distribution
FROM kuesioner_materialized k
WHERE k.Judul IS NOT NULL AND k.Judul != '' AND (k.KodeFakultas IS NULL OR k.KodeFakultas = '')
GROUP BY k.Judul, k.Semester, COALESCE(NULLIF(k.Unit, ''), 'Umum'), COALESCE(k.Fakultas, '');
""")

# c) Level 3 Per Prodi Filter
run_sql("""
INSERT INTO report_distribusi_fakultas (judul, semester, kode_fakultas, kode_prodi, unit, nama_fakultas, total_responden, persentase, prodi_distribution)
SELECT 
    k.Judul AS judul,
    CAST(k.Semester AS CHAR) AS semester,
    k.KodeFakultas AS kode_fakultas,
    k.KodeProdi AS kode_prodi,
    '' AS unit,
    COALESCE(k.Fakultas, '') AS nama_fakultas,
    COUNT(DISTINCT COALESCE(NULLIF(k.NPM,''), NULLIF(k.NIDN,''), NULLIF(k.NIP,''), k.NamaMahasiswa, k.NamaDosen, k.NamaTendik)) AS total_responden,
    100.00 AS persentase,
    CONCAT('[{"title":"', REPLACE(COALESCE(NULLIF(k.Prodi, ''), 'Umum'), '"', '\\\\\"'), '","total":', COUNT(DISTINCT COALESCE(NULLIF(k.NPM,''), NULLIF(k.NIDN,''), NULLIF(k.NIP,''), k.NamaMahasiswa, k.NamaDosen, k.NamaTendik)), '}]') AS prodi_distribution
FROM kuesioner_materialized k
WHERE k.Judul IS NOT NULL AND k.Judul != '' AND k.KodeFakultas IS NOT NULL AND k.KodeFakultas != '' AND k.KodeProdi IS NOT NULL AND k.KodeProdi != ''
GROUP BY k.Judul, k.Semester, k.KodeFakultas, k.KodeProdi, k.Fakultas;
""")

# Update percentages for Level 1
run_sql("""
UPDATE report_distribusi_fakultas d
JOIN (
    SELECT Judul, Semester, COUNT(DISTINCT COALESCE(NULLIF(NPM,''), NULLIF(NIDN,''), NULLIF(NIP,''), NamaMahasiswa, NamaDosen, NamaTendik)) AS grand_total
    FROM kuesioner_materialized
    WHERE Judul IS NOT NULL AND Judul != ''
    GROUP BY Judul, Semester
) tot ON tot.Judul = d.judul AND tot.Semester = d.semester
SET d.persentase = LEAST(100.00, ROUND((d.total_responden / NULLIF(tot.grand_total, 0)) * 100, 2))
WHERE d.kode_prodi = '';
""")
print("2. report_distribusi_fakultas populated!")

# 3. Populate report_top_questions (Overall, Per Fakultas, Per Prodi, Per Unit)
# a) Overall
run_sql("""
INSERT INTO report_top_questions (judul, semester, kode_fakultas, kode_prodi, unit, pertanyaan, nama_kategori, total_responden, rata_rata_skor, peringkat)
SELECT judul, semester, kode_fakultas, kode_prodi, unit, pertanyaan, nama_kategori, total_responden, rata_rata_skor, peringkat
FROM (
    SELECT 
        k.Judul AS judul,
        CAST(k.Semester AS CHAR) AS semester,
        '' AS kode_fakultas,
        '' AS kode_prodi,
        '' AS unit,
        k.Pertanyaan AS pertanyaan,
        COALESCE(NULLIF(k.Kategori, ''), 'Umum') AS nama_kategori,
        COUNT(DISTINCT COALESCE(NULLIF(k.NPM,''), NULLIF(k.NIDN,''), NULLIF(k.NIP,''), k.NamaMahasiswa, k.NamaDosen, k.NamaTendik)) AS total_responden,
        ROUND(AVG(CASE WHEN k.JenisPilihan = 'rating' AND k.Jawaban REGEXP '^[0-9]+$' THEN CAST(k.Jawaban AS DECIMAL(4,2)) ELSE NULL END), 2) AS rata_rata_skor,
        ROW_NUMBER() OVER (
            PARTITION BY k.Judul, k.Semester 
            ORDER BY AVG(CASE WHEN k.JenisPilihan = 'rating' AND k.Jawaban REGEXP '^[0-9]+$' THEN CAST(k.Jawaban AS DECIMAL(4,2)) ELSE NULL END) DESC
        ) AS peringkat
    FROM kuesioner_materialized k
    WHERE k.Judul IS NOT NULL AND k.Judul != '' AND k.Pertanyaan IS NOT NULL AND k.Pertanyaan != ''
    GROUP BY k.Judul, k.Semester, k.Pertanyaan, COALESCE(NULLIF(k.Kategori, ''), 'Umum')
) top_q
WHERE peringkat <= 10;
""")

# b) Per Fakultas
run_sql("""
INSERT INTO report_top_questions (judul, semester, kode_fakultas, kode_prodi, unit, pertanyaan, nama_kategori, total_responden, rata_rata_skor, peringkat)
SELECT judul, semester, kode_fakultas, kode_prodi, unit, pertanyaan, nama_kategori, total_responden, rata_rata_skor, peringkat
FROM (
    SELECT 
        k.Judul AS judul,
        CAST(k.Semester AS CHAR) AS semester,
        k.KodeFakultas AS kode_fakultas,
        '' AS kode_prodi,
        '' AS unit,
        k.Pertanyaan AS pertanyaan,
        COALESCE(NULLIF(k.Kategori, ''), 'Umum') AS nama_kategori,
        COUNT(DISTINCT COALESCE(NULLIF(k.NPM,''), NULLIF(k.NIDN,''), NULLIF(k.NIP,''), k.NamaMahasiswa, k.NamaDosen, k.NamaTendik)) AS total_responden,
        ROUND(AVG(CASE WHEN k.JenisPilihan = 'rating' AND k.Jawaban REGEXP '^[0-9]+$' THEN CAST(k.Jawaban AS DECIMAL(4,2)) ELSE NULL END), 2) AS rata_rata_skor,
        ROW_NUMBER() OVER (
            PARTITION BY k.Judul, k.Semester, k.KodeFakultas
            ORDER BY AVG(CASE WHEN k.JenisPilihan = 'rating' AND k.Jawaban REGEXP '^[0-9]+$' THEN CAST(k.Jawaban AS DECIMAL(4,2)) ELSE NULL END) DESC
        ) AS peringkat
    FROM kuesioner_materialized k
    WHERE k.Judul IS NOT NULL AND k.Judul != '' AND k.Pertanyaan IS NOT NULL AND k.Pertanyaan != '' AND k.KodeFakultas IS NOT NULL AND k.KodeFakultas != ''
    GROUP BY k.Judul, k.Semester, k.KodeFakultas, k.Pertanyaan, COALESCE(NULLIF(k.Kategori, ''), 'Umum')
) top_q
WHERE peringkat <= 10;
""")

# c) Per Prodi
run_sql("""
INSERT INTO report_top_questions (judul, semester, kode_fakultas, kode_prodi, unit, pertanyaan, nama_kategori, total_responden, rata_rata_skor, peringkat)
SELECT judul, semester, kode_fakultas, kode_prodi, unit, pertanyaan, nama_kategori, total_responden, rata_rata_skor, peringkat
FROM (
    SELECT 
        k.Judul AS judul,
        CAST(k.Semester AS CHAR) AS semester,
        k.KodeFakultas AS kode_fakultas,
        k.KodeProdi AS kode_prodi,
        '' AS unit,
        k.Pertanyaan AS pertanyaan,
        COALESCE(NULLIF(k.Kategori, ''), 'Umum') AS nama_kategori,
        COUNT(DISTINCT COALESCE(NULLIF(k.NPM,''), NULLIF(k.NIDN,''), NULLIF(k.NIP,''), k.NamaMahasiswa, k.NamaDosen, k.NamaTendik)) AS total_responden,
        ROUND(AVG(CASE WHEN k.JenisPilihan = 'rating' AND k.Jawaban REGEXP '^[0-9]+$' THEN CAST(k.Jawaban AS DECIMAL(4,2)) ELSE NULL END), 2) AS rata_rata_skor,
        ROW_NUMBER() OVER (
            PARTITION BY k.Judul, k.Semester, k.KodeFakultas, k.KodeProdi
            ORDER BY AVG(CASE WHEN k.JenisPilihan = 'rating' AND k.Jawaban REGEXP '^[0-9]+$' THEN CAST(k.Jawaban AS DECIMAL(4,2)) ELSE NULL END) DESC
        ) AS peringkat
    FROM kuesioner_materialized k
    WHERE k.Judul IS NOT NULL AND k.Judul != '' AND k.Pertanyaan IS NOT NULL AND k.Pertanyaan != '' AND k.KodeProdi IS NOT NULL AND k.KodeProdi != ''
    GROUP BY k.Judul, k.Semester, k.KodeFakultas, k.KodeProdi, k.Pertanyaan, COALESCE(NULLIF(k.Kategori, ''), 'Umum')
) top_q
WHERE peringkat <= 10;
""")

# d) Per Unit
run_sql("""
INSERT INTO report_top_questions (judul, semester, kode_fakultas, kode_prodi, unit, pertanyaan, nama_kategori, total_responden, rata_rata_skor, peringkat)
SELECT judul, semester, kode_fakultas, kode_prodi, unit, pertanyaan, nama_kategori, total_responden, rata_rata_skor, peringkat
FROM (
    SELECT 
        k.Judul AS judul,
        CAST(k.Semester AS CHAR) AS semester,
        '' AS kode_fakultas,
        '' AS kode_prodi,
        COALESCE(NULLIF(k.Unit, ''), 'Umum') AS unit,
        k.Pertanyaan AS pertanyaan,
        COALESCE(NULLIF(k.Kategori, ''), 'Umum') AS nama_kategori,
        COUNT(DISTINCT COALESCE(NULLIF(k.NPM,''), NULLIF(k.NIDN,''), NULLIF(k.NIP,''), k.NamaMahasiswa, k.NamaDosen, k.NamaTendik)) AS total_responden,
        ROUND(AVG(CASE WHEN k.JenisPilihan = 'rating' AND k.Jawaban REGEXP '^[0-9]+$' THEN CAST(k.Jawaban AS DECIMAL(4,2)) ELSE NULL END), 2) AS rata_rata_skor,
        ROW_NUMBER() OVER (
            PARTITION BY k.Judul, k.Semester, COALESCE(NULLIF(k.Unit, ''), 'Umum')
            ORDER BY AVG(CASE WHEN k.JenisPilihan = 'rating' AND k.Jawaban REGEXP '^[0-9]+$' THEN CAST(k.Jawaban AS DECIMAL(4,2)) ELSE NULL END) DESC
        ) AS peringkat
    FROM kuesioner_materialized k
    WHERE k.Judul IS NOT NULL AND k.Judul != '' AND k.Pertanyaan IS NOT NULL AND k.Pertanyaan != '' AND (k.KodeFakultas IS NULL OR k.KodeFakultas = '')
    GROUP BY k.Judul, k.Semester, COALESCE(NULLIF(k.Unit, ''), 'Umum'), k.Pertanyaan, COALESCE(NULLIF(k.Kategori, ''), 'Umum')
) top_q
WHERE peringkat <= 10;
""")
print("3. report_top_questions populated!")

# 4. Populate report_kategori_summary using Python aggregation for perfect questions_json
q_raw = run_sql("""
    SELECT 
        COALESCE(NULLIF(KodeFakultas, ''), '') AS k_fak,
        COALESCE(NULLIF(KodeProdi, ''), '') AS k_prodi,
        COALESCE(NULLIF(Unit, ''), 'Umum') AS unit,
        Judul AS judul,
        CAST(Semester AS CHAR) AS semester,
        COALESCE(NULLIF(Kategori, ''), 'Umum') AS nama_kategori,
        MAX(FullPath) AS full_text,
        Pertanyaan AS pertanyaan,
        COALESCE(NULLIF(JenisPilihan, ''), 'rating') AS jenispilihan,
        SUM(CASE WHEN Jawaban = '1' OR Jawaban = 'Sangat Tidak Baik' OR Jawaban = 'Sangat Tidak Setuju' THEN 1 ELSE 0 END) AS r1,
        SUM(CASE WHEN Jawaban = '2' OR Jawaban = 'Tidak Baik' OR Jawaban = 'Tidak Setuju' THEN 1 ELSE 0 END) AS r2,
        SUM(CASE WHEN Jawaban = '3' OR Jawaban = 'Cukup' OR Jawaban = 'Netral' THEN 1 ELSE 0 END) AS r3,
        SUM(CASE WHEN Jawaban = '4' OR Jawaban = 'Baik' OR Jawaban = 'Setuju' THEN 1 ELSE 0 END) AS r4,
        SUM(CASE WHEN Jawaban = '5' OR Jawaban = 'Sangat Baik' OR Jawaban = 'Sangat Setuju' THEN 1 ELSE 0 END) AS r5,
        COUNT(*) AS total_resp,
        COALESCE(AVG(
            CASE 
                WHEN Jawaban = '1' OR Jawaban = 'Sangat Tidak Baik' OR Jawaban = 'Sangat Tidak Setuju' THEN 1.0
                WHEN Jawaban = '2' OR Jawaban = 'Tidak Baik' OR Jawaban = 'Tidak Setuju' THEN 2.0
                WHEN Jawaban = '3' OR Jawaban = 'Cukup' OR Jawaban = 'Netral' THEN 3.0
                WHEN Jawaban = '4' OR Jawaban = 'Baik' OR Jawaban = 'Setuju' THEN 4.0
                WHEN Jawaban = '5' OR Jawaban = 'Sangat Baik' OR Jawaban = 'Sangat Setuju' THEN 5.0
                WHEN Jawaban REGEXP '^[0-9]+(\\.[0-9]+)?$' THEN CAST(Jawaban AS DECIMAL(4,2))
                ELSE NULL 
            END
        ), 0.0) AS avg_score
    FROM kuesioner_materialized
    WHERE Judul IS NOT NULL AND Judul != '' AND Kategori IS NOT NULL AND Kategori != '' AND Pertanyaan IS NOT NULL AND Pertanyaan != ''
    GROUP BY KodeFakultas, KodeProdi, COALESCE(NULLIF(Unit, ''), 'Umum'), Judul, Semester, COALESCE(NULLIF(Kategori, ''), 'Umum'), Pertanyaan, COALESCE(NULLIF(JenisPilihan, ''), 'rating')
""")

overall_map = {}
fak_map = {}
prodi_map = {}
unit_map = {}

for line in q_raw.split("\n"):
    if not line.strip():
        continue
    parts = line.split("\t")
    if len(parts) < 16:
        continue
    k_fak, k_prodi, unit_val, judul, semester, kat, full_text, pert, j_pilihan = parts[0], parts[1], parts[2], parts[3], parts[4], parts[5], parts[6], parts[7], parts[8]
    r1, r2, r3, r4, r5 = int(parts[9]), int(parts[10]), int(parts[11]), int(parts[12]), int(parts[13])
    t_resp = int(parts[14])
    avg_s = float(parts[15]) if parts[15] else 0.0

    def add_to_map(target_dict, key):
        if key not in target_dict:
            target_dict[key] = {"full_text": full_text, "questions": {}}
        if pert not in target_dict[key]["questions"]:
            target_dict[key]["questions"][pert] = {"title": pert, "jenispilihan": j_pilihan, "chart_distribution": {"1": 0, "2": 0, "3": 0, "4": 0, "5": 0}, "total_resp": 0, "scores": []}
        q = target_dict[key]["questions"][pert]
        q["chart_distribution"]["1"] += r1
        q["chart_distribution"]["2"] += r2
        q["chart_distribution"]["3"] += r3
        q["chart_distribution"]["4"] += r4
        q["chart_distribution"]["5"] += r5
        q["total_resp"] += t_resp
        if avg_s > 0: q["scores"].append(avg_s)

    # 1. Overall
    add_to_map(overall_map, (judul, semester, "", "", "", kat))
    
    # 2. Per Fakultas (only if k_fak exists)
    if k_fak:
        add_to_map(fak_map, (judul, semester, k_fak, "", "", kat))
        
    # 3. Per Prodi (only if k_prodi exists)
    if k_fak and k_prodi:
        add_to_map(prodi_map, (judul, semester, k_fak, k_prodi, "", kat))
        
    # 4. Per Unit (only if k_fak is empty)
    if not k_fak:
        add_to_map(unit_map, (judul, semester, "", "", unit_val, kat))

combined_map = {}
combined_map.update(overall_map)
combined_map.update(fak_map)
combined_map.update(prodi_map)
combined_map.update(unit_map)

values_list = []
for (judul, sem, k_fak, k_prodi, unit_val, kat), data in combined_map.items():
    q_arr = []
    c_chart = {"1": 0, "2": 0, "3": 0, "4": 0, "5": 0}
    scores = []
    tot_resp = 0
    tot_q = len(data["questions"])
    
    for pert, q in data["questions"].items():
        q_arr.append({
            "title": q["title"],
            "jenispilihan": q["jenispilihan"],
            "chart_distribution": q["chart_distribution"]
        })
        for r_k in ["1", "2", "3", "4", "5"]:
            c_chart[r_k] += q["chart_distribution"][r_k]
        tot_resp += q["total_resp"]
        if q["scores"]:
            scores.extend(q["scores"])
            
    avg_score = round(sum(scores) / len(scores), 2) if scores else 0.0
    q_json = json.dumps(q_arr)
    c_json = json.dumps(c_chart)
    
    values_list.append(f"({escape_sql(judul)}, {escape_sql(sem)}, {escape_sql(k_fak)}, {escape_sql(k_prodi)}, {escape_sql(unit_val)}, {escape_sql(kat)}, {escape_sql(data['full_text'])}, {tot_q}, {tot_resp}, {avg_score}, {escape_sql(c_json)}, {escape_sql(q_json)})")

if values_list:
    for i in range(0, len(values_list), 50):
        chunk = values_list[i:i+50]
        run_sql(f"INSERT INTO report_kategori_summary (judul, semester, kode_fakultas, kode_prodi, unit, nama_kategori, full_text, total_pertanyaan, total_responden, rata_rata_skor, chart_distribution, questions_json) VALUES {','.join(chunk)}")

print("4. report_kategori_summary populated!")

# 5. Populate report_year
run_sql("""
INSERT INTO report_year (tahun, total_kuesioner, total_mahasiswa, total_dosen, total_tendik)
SELECT 
    LEFT(CAST(Semester AS CHAR), 4) AS tahun,
    COUNT(*) AS total_kuesioner,
    COUNT(DISTINCT CASE WHEN NPM != '' OR NamaMahasiswa != '' THEN COALESCE(NULLIF(NPM,''), NamaMahasiswa) END) AS total_mahasiswa,
    COUNT(DISTINCT CASE WHEN NIDN != '' OR NamaDosen != '' THEN COALESCE(NULLIF(NIDN,''), NamaDosen) END) AS total_dosen,
    COUNT(DISTINCT CASE WHEN NIP != '' OR NamaTendik != '' THEN COALESCE(NULLIF(NIP,''), NamaTendik) END) AS total_tendik
FROM kuesioner_materialized
WHERE Semester IS NOT NULL AND Semester != ''
GROUP BY LEFT(CAST(Semester AS CHAR), 4)
ORDER BY tahun ASC;
""")
print("5. report_year populated!")
print("=== ALL 5 TABLES RECALCULATED INSTANTLY WITH PERFECT FAKULTAS, PRODI & UNIT DATA! ===")
PYTHON_SCRIPT

EXIT_CODE=$?

if [ $EXIT_CODE -eq 0 ]; then
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Kalkulasi pure bash report summary 5 tabel berhasil! 🚀" >> "$LOG_FILE"
    echo "Kalkulasi pure bash report summary 5 tabel berhasil! 🚀"
else
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Kalkulasi GAGAL" >> "$ERROR_LOG"
    echo "Kalkulasi GAGAL!"
    exit 1
fi
