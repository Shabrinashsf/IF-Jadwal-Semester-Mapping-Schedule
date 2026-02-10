# Excel Reader Go

Aplikasi sederhana untuk membaca data dari file Excel (.xlsx) menggunakan Go.

## Requirements

- Go 1.21+

## Installation

```bash
go mod tidy
```

## Usage

1. Letakkan file Excel kamu (misal `data.xlsx`) di folder project
2. Jalankan:

```bash
go run main.go
```

## Fitur

- Membaca semua baris dan kolom dari sheet
- Membaca cell spesifik (misal A1, B2)
- Membaca range tertentu (misal A1:C3)

## Contoh Output

```
Membaca sheet: Sheet1

=== Semua Data ===
Baris 1: [Nama Umur Kota]
Baris 2: [Budi 25 Jakarta]
Baris 3: [Ani 30 Surabaya]

=== Baca Cell Spesifik ===
A1: Nama
B2: 25

=== Baca Range A1:C3 ===
A1: Nama       B1: Umur       C1: Kota       
A2: Budi       B2: 25         C2: Jakarta    
A3: Ani        B3: 30         C3: Surabaya   
```

## Library

- [excelize](https://github.com/xuri/excelize) - Library Go untuk membaca/menulis Excel
