# IF Jadwal Semester - Mapping Schedule

Aplikasi untuk membaca dan memetakan jadwal kuliah dari file Excel, kemudian mengorganisir dan mengekspor data berdasarkan program studi, semester, dan informasi mata kuliah lainnya.

## Deskripsi

Aplikasi ini membaca file Excel jadwal kuliah (format: `data.xlsx`) dan melakukan:
- **Parsing**: Mengekstrak informasi jadwal kuliah (hari, jam, ruangan, mata kuliah, semester, dosen, SKS)
- **Mapping**: Mengelompokkan jadwal berdasarkan semester dan program studi
- **Validasi**: Memastikan format data sesuai dengan pola yang ditentukan
- **Ekspor**: Menghasilkan output dalam format JSON dan Excel yang terstruktur

## Requirements

- Go 1.21+
- Library: [excelize](https://github.com/xuri/excelize) v2

## Installation

```bash
go mod tidy
```

## Cara Penggunaan

1. Siapkan file Excel `data.xlsx` dengan struktur:
   - **Baris 1**: Header dengan nama-nama ruangan (misal: IF-101, IF-102, dll)
   - **Kolom A**: Hari (SENIN, SELASA, RABU, KAMIS, JUMAT)
   - **Kolom B**: Jam (misal: 07.00 - 07.50)
   - **Mulai dari Kolom G**: Data mata kuliah dengan format `PRODI_NamaMataKuliah` (baris berikutnya: `Sem X / KODE_DOSEN / SKS`)

2. Letakkan file `data.xlsx` di folder project

3. Jalankan aplikasi:
```bash
go run main.go
```

## Output

Aplikasi menghasilkan dua file output:

### 1. jadwal.json
File JSON yang berisi data terstruktur dengan format `{prodi: {semester: [jadwal]}}`
- Contoh struktur: `{"IF": {"1": [...], "2": [...]}, "S2": {"1": [...]}}`
- Jadwal dikelompokkan per program studi dan semester

### 2. jadwal.xlsx
File Excel dengan multiple sheet, masing-masing sheet berisi:
- Format nama sheet: `{PRODI}_Sem_{Semester}` (misal: `IF_Sem_1`, `RKA_Sem_3`)
- Kolom: Hari, Jam, Ruangan, Prodi, Mata Kuliah, Semester, Kode Dosen, SKS

## Fitur Utama

- ✓ Membaca struktur jadwal kompleks dari Excel
- ✓ Ekstraksi informasi jadwal per mata kuliah (hari, jam, ruangan, dosen, SKS)
- ✓ Parsing program studi dari nama mata kuliah (IF, IUP, RKA, RPL)
- ✓ Pembersihan data ruangan (menghapus suffix seperti "a&b")
- ✓ Pengelompokan otomatis berdasarkan semester dan prodi
- ✓ Ekspor ke JSON dengan struktur yang terorganisir
- ✓ Ekspor ke Excel dengan sheet terpisah per prodi/semester

## Contoh Output Console

```
Membaca sheet: Jadwal Kuliah

=== SEMUA JADWAL ===
Total: 24 jadwal ditemukan

=== JADWAL PER PRODI & SEMESTER ===
  IF:
    Semester 1: 6 mata kuliah
    Semester 3: 5 mata kuliah
  RKA:
    Semester 1: 4 mata kuliah
    Semester 3: 4 mata kuliah

✓ Berhasil export ke jadwal.json

=== EXPORT KE EXCEL ===
Export jadwal per prodi dan semester...
  ✓ Sheet 'IF_Sem_1' (6 data)
  ✓ Sheet 'IF_Sem_3' (5 data)
  ✓ Sheet 'RKA_Sem_1' (4 data)
  ✓ Sheet 'RKA_Sem_3' (4 data)

✓ Berhasil export ke jadwal.xlsx
```

## Struktur Data

Setiap jadwal memiliki informasi:
```
Hari       : SENIN, SELASA, RABU, KAMIS, atau JUMAT
Jam        : HH.MM - HH.MM (contoh: 07.00 - 07.50)
Ruangan    : Kode ruangan (contoh: IF-101, IF-AV)
Prodi      : Program studi (IF, IUP, RKA, RPL)
Mata Kuliah: Nama mata kuliah
Semester   : 1-8
Kode Dosen : Inisial dosen (2 huruf)
SKS        : Jumlah satuan kredit semester
```

## Library Utama

- [excelize](https://github.com/xuri/excelize) - Library Go untuk membaca & menulis file Excel (.xlsx)
