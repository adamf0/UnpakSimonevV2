#!/bin/bash

# ==============================================================================
# SIMONEV - BATCH REPORT SUMMARY CALCULATOR SCRIPT (PURE SH / MYSQL CLI)
# ==============================================================================
# Script ini 100% pure .sh / python3 runner (tanpa kode Golang / .go file).
# Mengkalkulasi data dari kuesioner_materialized ke 5 tabel summary:
# 1. report_summary_overview    (Data Overview Utama)
# 2. report_distribusi_fakultas (Distribusi Responden per Fakultas + Prodi JSON)
# 3. report_top_questions       (Top 10 High-Engagement Questions)
# 4. report_kategori_summary     (Ringkasan Rating & JSON Chart/Pertanyaan per Kategori)
# 5. report_year                 (Grafik 4 Year Chart: Mahasiswa, Dosen, Tendik)
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
    cmd = [mysql_bin, f"-h{host}", f"-P{port}", f"-u{user}"]
    if password:
        cmd.append(f"-p{password}")
    cmd.extend([name, "-N", "-B", "-e", sql])
    
    p = subprocess.run(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    if p.returncode != 0:
        raise Exception(f"SQL Error: {p.stderr.strip()}")
    return p.stdout.strip()

def escape_sql(val):
    if val is None:
        return "NULL"
    val_str = str(val).replace("\\", "\\\\").replace("'", "''")
    return f"'{val_str}'"

print("=== STARTING PURE BASH CALCULATOR FOR 5 SUMMARY TABLES ===")

# Create tables if not exist
ddls = [
    "DROP TABLE IF EXISTS report_summary_overview, report_distribusi_fakultas, report_top_questions, report_kategori_summary, report_year",
    """CREATE TABLE IF NOT EXISTS report_summary_overview (
        id INT AUTO_INCREMENT PRIMARY KEY,
        judul VARCHAR(255) NOT NULL,
        semester VARCHAR(50) NOT NULL,
        total_responden INT DEFAULT 0,
        total_jawaban INT DEFAULT 0,
        rata_rata_rating DECIMAL(4,2) DEFAULT 0.00,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        UNIQUE KEY uq_judul_sem (judul(150), semester)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4""",
    """CREATE TABLE IF NOT EXISTS report_distribusi_fakultas (
        id INT AUTO_INCREMENT PRIMARY KEY,
        judul VARCHAR(255) NOT NULL,
        semester VARCHAR(50) NOT NULL,
        kode_fakultas VARCHAR(50) NOT NULL,
        nama_fakultas VARCHAR(255),
        total_responden INT DEFAULT 0,
        persentase DECIMAL(5,2) DEFAULT 0.00,
        prodi_distribution JSON,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        UNIQUE KEY uq_judul_sem_fak (judul(150), semester, kode_fakultas)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4""",
    """CREATE TABLE IF NOT EXISTS report_top_questions (
        id INT AUTO_INCREMENT PRIMARY KEY,
        judul VARCHAR(255) NOT NULL,
        semester VARCHAR(50) NOT NULL,
        pertanyaan TEXT NOT NULL,
        nama_kategori VARCHAR(255),
        total_responden INT DEFAULT 0,
        rata_rata_skor DECIMAL(4,2) DEFAULT 0.00,
        peringkat INT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        KEY idx_judul_sem_rank (judul(150), semester, peringkat)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4""",
    """CREATE TABLE IF NOT EXISTS report_kategori_summary (
        id INT AUTO_INCREMENT PRIMARY KEY,
        judul VARCHAR(255) NOT NULL,
        semester VARCHAR(50) NOT NULL,
        nama_kategori VARCHAR(255) NOT NULL,
        full_text TEXT,
        total_pertanyaan INT DEFAULT 0,
        total_responden INT DEFAULT 0,
        rata_rata_skor DECIMAL(4,2) DEFAULT 0.00,
        chart_distribution JSON,
        questions_json JSON,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        KEY idx_judul_sem_kat (judul(150), semester, nama_kategori(100))
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
run_sql("""
INSERT INTO report_summary_overview (judul, semester, total_responden, total_jawaban, rata_rata_rating)
SELECT 
    k.Judul AS judul,
    CAST(k.Semester AS CHAR) AS semester,
    COUNT(DISTINCT COALESCE(NULLIF(k.NPM,''), NULLIF(k.NIDN,''), NULLIF(k.NIP,''), k.NamaMahasiswa, k.NamaDosen, k.NamaTendik)) AS total_responden,
    COUNT(*) AS total_jawaban,
    COALESCE(ROUND(AVG(CASE WHEN k.JenisPilihan = 'rating' AND k.Jawaban REGEXP '^[0-9]+$' THEN CAST(k.Jawaban AS DECIMAL(4,2)) ELSE NULL END), 2), 0.00) AS rata_rata_rating
FROM kuesioner_materialized k
WHERE k.Judul IS NOT NULL AND k.Judul != ''
GROUP BY k.Judul, k.Semester;
""")
print("1. report_summary_overview populated!")

# Get distinct Judul and Semester
js_raw = run_sql("SELECT DISTINCT Judul, CAST(Semester AS CHAR) FROM kuesioner_materialized WHERE Judul IS NOT NULL AND Judul != ''")

# 2. Populate report_distribusi_fakultas
for line in js_raw.split("\n"):
    if not line.strip():
        continue
    parts = line.split("\t")
    if len(parts) < 2:
        continue
    judul, semester = parts[0], parts[1]
    
    gt_str = run_sql(f"""
        SELECT COUNT(DISTINCT COALESCE(NULLIF(NPM,''), NULLIF(NIDN,''), NULLIF(NIP,''), NamaMahasiswa, NamaDosen, NamaTendik))
        FROM kuesioner_materialized
        WHERE Judul = {escape_sql(judul)} AND CAST(Semester AS CHAR) = {escape_sql(semester)}
    """)
    grand_total = int(gt_str) if gt_str and gt_str.isdigit() else 0
    if grand_total == 0:
        continue
    
    f_raw = run_sql(f"""
        SELECT 
            COALESCE(NULLIF(KodeFakultas, ''), 'UNKNOWN') AS k_fak,
            COALESCE(NULLIF(Fakultas, ''), 'Tidak Terdefinisi / Tendik') AS n_fak,
            COALESCE(NULLIF(Prodi, ''), 'Umum') AS n_prodi,
            COUNT(DISTINCT COALESCE(NULLIF(NPM,''), NULLIF(NIDN,''), NULLIF(NIP,''), NamaMahasiswa, NamaDosen, NamaTendik)) AS total_resp
        FROM kuesioner_materialized
        WHERE Judul = {escape_sql(judul)} AND CAST(Semester AS CHAR) = {escape_sql(semester)}
        GROUP BY COALESCE(NULLIF(KodeFakultas, ''), 'UNKNOWN'), COALESCE(NULLIF(Fakultas, ''), 'Tidak Terdefinisi / Tendik'), COALESCE(NULLIF(Prodi, ''), 'Umum')
        ORDER BY MIN(id) ASC
    """)
    
    fak_map = {}
    fak_order = []
    for f_line in f_raw.split("\n"):
        if not f_line.strip():
            continue
        f_parts = f_line.split("\t")
        if len(f_parts) < 4:
            continue
        k_fak, n_fak, n_prodi, count_str = f_parts[0], f_parts[1], f_parts[2], f_parts[3]
        count = int(count_str) if count_str.isdigit() else 0
        
        if k_fak not in fak_map:
            fak_map[k_fak] = {"nama": n_fak, "total": 0, "prodi": []}
            fak_order.append(k_fak)
        fak_map[k_fak]["total"] += count
        fak_map[k_fak]["prodi"].append({"title": n_prodi, "total": count})
        
    for k_fak in fak_order:
        data = fak_map[k_fak]
        pct = round((data["total"] / grand_total) * 100, 2)
        prodi_json = json.dumps(data["prodi"])
        
        run_sql(f"""
            INSERT INTO report_distribusi_fakultas (judul, semester, kode_fakultas, nama_fakultas, total_responden, persentase, prodi_distribution)
            VALUES ({escape_sql(judul)}, {escape_sql(semester)}, {escape_sql(k_fak)}, {escape_sql(data["nama"])}, {data["total"]}, {pct}, {escape_sql(prodi_json)})
        """)
print("2. report_distribusi_fakultas populated!")

# 3. Populate report_top_questions
for line in js_raw.split("\n"):
    if not line.strip():
        continue
    parts = line.split("\t")
    if len(parts) < 2:
        continue
    judul, semester = parts[0], parts[1]
    
    top_raw = run_sql(f"""
        SELECT 
            Pertanyaan,
            COALESCE(NULLIF(Kategori, ''), 'Umum') AS kat,
            COUNT(DISTINCT COALESCE(NULLIF(NPM,''), NULLIF(NIDN,''), NULLIF(NIP,''), NamaMahasiswa, NamaDosen, NamaTendik)) AS total_resp,
            ROUND(AVG(CASE WHEN JenisPilihan = 'rating' AND Jawaban REGEXP '^[0-9]+$' THEN CAST(Jawaban AS DECIMAL(4,2)) ELSE NULL END), 2) AS avg_skor
        FROM kuesioner_materialized
        WHERE Judul = {escape_sql(judul)} AND CAST(Semester AS CHAR) = {escape_sql(semester)} AND Pertanyaan IS NOT NULL AND Pertanyaan != ''
        GROUP BY Pertanyaan, COALESCE(NULLIF(Kategori, ''), 'Umum')
        ORDER BY avg_skor DESC, MIN(id) ASC
        LIMIT 10
    """)
    
    rank = 1
    for t_line in top_raw.split("\n"):
        if not t_line.strip():
            continue
        t_parts = t_line.split("\t")
        if len(t_parts) < 4:
            continue
        q_title, q_kat, t_resp, avg_skor = t_parts[0], t_parts[1], t_parts[2], t_parts[3]
        avg_val = float(avg_skor) if avg_skor else 0.0
        
        run_sql(f"""
            INSERT INTO report_top_questions (judul, semester, pertanyaan, nama_kategori, total_responden, rata_rata_skor, peringkat)
            VALUES ({escape_sql(judul)}, {escape_sql(semester)}, {escape_sql(q_title)}, {escape_sql(q_kat)}, {t_resp}, {avg_val}, {rank})
        """)
        rank += 1
print("3. report_top_questions populated!")

# 4. Populate report_kategori_summary
for line in js_raw.split("\n"):
    if not line.strip():
        continue
    parts = line.split("\t")
    if len(parts) < 2:
        continue
    judul, semester = parts[0], parts[1]
    
    cat_raw = run_sql(f"""
        SELECT 
            COALESCE(NULLIF(Kategori, ''), 'Umum') AS kat,
            MAX(FullPath) AS full_p,
            COUNT(DISTINCT Pertanyaan) AS total_q,
            COUNT(*) AS total_resp,
            ROUND(AVG(CASE WHEN JenisPilihan = 'rating' AND Jawaban REGEXP '^[0-9]+$' THEN CAST(Jawaban AS DECIMAL(4,2)) ELSE NULL END), 2) AS avg_skor,
            SUM(CASE WHEN Jawaban = '1' OR Jawaban = 'Sangat Tidak Baik' OR Jawaban = 'Sangat Tidak Setuju' THEN 1 ELSE 0 END) AS r1,
            SUM(CASE WHEN Jawaban = '2' OR Jawaban = 'Tidak Baik' OR Jawaban = 'Tidak Setuju' THEN 1 ELSE 0 END) AS r2,
            SUM(CASE WHEN Jawaban = '3' OR Jawaban = 'Cukup' OR Jawaban = 'Netral' THEN 1 ELSE 0 END) AS r3,
            SUM(CASE WHEN Jawaban = '4' OR Jawaban = 'Baik' OR Jawaban = 'Setuju' THEN 1 ELSE 0 END) AS r4,
            SUM(CASE WHEN Jawaban = '5' OR Jawaban = 'Sangat Baik' OR Jawaban = 'Sangat Setuju' THEN 1 ELSE 0 END) AS r5
        FROM kuesioner_materialized
        WHERE Judul = {escape_sql(judul)} AND CAST(Semester AS CHAR) = {escape_sql(semester)} AND Kategori IS NOT NULL AND Kategori != ''
        GROUP BY COALESCE(NULLIF(Kategori, ''), 'Umum')
        ORDER BY MIN(id) ASC
    """)
    
    for c_line in cat_raw.split("\n"):
        if not c_line.strip():
            continue
        c_parts = c_line.split("\t")
        if len(c_parts) < 10:
            continue
        kat, full_p, total_q, total_resp, avg_skor = c_parts[0], c_parts[1], c_parts[2], c_parts[3], c_parts[4]
        r1, r2, r3, r4, r5 = c_parts[5], c_parts[6], c_parts[7], c_parts[8], c_parts[9]
        
        chart_dist = {
            "1": int(r1) if r1.isdigit() else 0,
            "2": int(r2) if r2.isdigit() else 0,
            "3": int(r3) if r3.isdigit() else 0,
            "4": int(r4) if r4.isdigit() else 0,
            "5": int(r5) if r5.isdigit() else 0,
        }
        
        q_raw = run_sql(f"""
            SELECT 
                Pertanyaan,
                JenisPilihan,
                SUM(CASE WHEN Jawaban = '1' OR Jawaban = 'Sangat Tidak Baik' OR Jawaban = 'Sangat Tidak Setuju' THEN 1 ELSE 0 END) AS r1,
                SUM(CASE WHEN Jawaban = '2' OR Jawaban = 'Tidak Baik' OR Jawaban = 'Tidak Setuju' THEN 1 ELSE 0 END) AS r2,
                SUM(CASE WHEN Jawaban = '3' OR Jawaban = 'Cukup' OR Jawaban = 'Netral' THEN 1 ELSE 0 END) AS r3,
                SUM(CASE WHEN Jawaban = '4' OR Jawaban = 'Baik' OR Jawaban = 'Setuju' THEN 1 ELSE 0 END) AS r4,
                SUM(CASE WHEN Jawaban = '5' OR Jawaban = 'Sangat Baik' OR Jawaban = 'Sangat Setuju' THEN 1 ELSE 0 END) AS r5
            FROM kuesioner_materialized
            WHERE Judul = {escape_sql(judul)} AND CAST(Semester AS CHAR) = {escape_sql(semester)} AND Kategori = {escape_sql(kat)} AND Pertanyaan != ''
            GROUP BY Pertanyaan, JenisPilihan
            ORDER BY MIN(id) ASC
        """)
        
        q_dists = []
        for q_line in q_raw.split("\n"):
            if not q_line.strip():
                continue
            q_parts = q_line.split("\t")
            if len(q_parts) < 7:
                continue
            q_title, q_jenis = q_parts[0], q_parts[1]
            qr1, qr2, qr3, qr4, qr5 = q_parts[2], q_parts[3], q_parts[4], q_parts[5], q_parts[6]
            q_dists.append({
                "title": q_title,
                "jenispilihan": q_jenis,
                "chart_distribution": {
                    "1": int(qr1) if qr1.isdigit() else 0,
                    "2": int(qr2) if qr2.isdigit() else 0,
                    "3": int(qr3) if qr3.isdigit() else 0,
                    "4": int(qr4) if qr4.isdigit() else 0,
                    "5": int(qr5) if qr5.isdigit() else 0,
                }
            })
            
        chart_json = json.dumps(chart_dist)
        questions_json = json.dumps(q_dists)
        avg_val = float(avg_skor) if avg_skor else 0.0
        
        run_sql(f"""
            INSERT INTO report_kategori_summary (judul, semester, nama_kategori, full_text, total_pertanyaan, total_responden, rata_rata_skor, chart_distribution, questions_json)
            VALUES ({escape_sql(judul)}, {escape_sql(semester)}, {escape_sql(kat)}, {escape_sql(full_p)}, {total_q}, {total_resp}, {avg_val}, {escape_sql(chart_json)}, {escape_sql(questions_json)})
        """)
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
print("=== ALL 5 TABLES RECALCULATED SUCCESSFULLY! ===")
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
