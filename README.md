# IF Jadwal Semester - Mapping Schedule

Application to read and map course schedules from Excel files, then organize and export data based on study program, semester, and other course information.

## Description

This application reads an Excel file containing course schedules (format: `data.xlsx`) and performs:
- **Parsing**: Extracts course schedule information (day, time, room, course, semester, lecturer, credits)
- **Mapping**: Groups schedules by semester and study program
- **Validation**: Ensures data format matches the specified patterns
- **Export**: Generates structured output in JSON and Excel formats

## Requirements

- Go 1.21+
- Library: [excelize](https://github.com/xuri/excelize) v2

## Installation

```bash
go mod tidy
```

## Usage

1. Prepare an Excel file `data.xlsx` with the following structure:
   - **Row 1**: Header with room names (e.g., IF-101, IF-102, etc.)
   - **Column A**: Day (SENIN, SELASA, RABU, KAMIS, JUMAT)
   - **Column B**: Time (e.g., 07.00 - 07.50)
   - **Starting from Column G**: Course data with format `PRODI_CourseName` (next row: `Sem X / LECTURER_CODE / SKS`)

2. Place the `data.xlsx` file in the project folder

3. Run the application:
```bash
go run main.go
```

## Output

The application generates two output files:

### 1. jadwal.json
JSON file containing structured data with format `{prodi: {semester: [schedule]}}`
- Example structure: `{"IF": {"1": [...], "2": [...]}, "S2": {"1": [...]}}`
- Schedules are grouped by study program and semester

### 2. jadwal.xlsx
Excel file with multiple sheets, each sheet containing:
- Sheet name format: `{PRODI}_Sem_{Semester}` (e.g., `IF_Sem_1`, `RKA_Sem_3`)
- Columns: Hari (Day), Jam (Time), Ruangan (Room), Prodi (Program), Mata Kuliah (Course), Semester, Kode Dosen (Lecturer Code), SKS (Credits)

## Key Features

- ✓ Reads complex schedule structures from Excel
- ✓ Extracts schedule information per course (day, time, room, lecturer, credits)
- ✓ Parses study program from course names (IF, IUP, RKA, RPL)
- ✓ Cleans room data (removes suffixes like "a&b")
- ✓ Automatic grouping by semester and study program
- ✓ Export to JSON with organized structure
- ✓ Export to Excel with separate sheets per program/semester

## Console Output Example

```
Reading sheet: Jadwal Kuliah

=== ALL SCHEDULES ===
Total: 24 schedules found

=== SCHEDULES BY PROGRAM & SEMESTER ===
  IF:
    Semester 1: 6 courses
    Semester 3: 5 courses
  RKA:
    Semester 1: 4 courses
    Semester 3: 4 courses

✓ Successfully exported to jadwal.json

=== EXPORT TO EXCEL ===
Exporting schedules per program and semester...
  ✓ Sheet 'IF_Sem_1' (6 data)
  ✓ Sheet 'IF_Sem_3' (5 data)
  ✓ Sheet 'RKA_Sem_1' (4 data)
  ✓ Sheet 'RKA_Sem_3' (4 data)

✓ Successfully exported to jadwal.xlsx
```

## Data Structure

Each schedule contains the following information:
```
Hari       : SENIN, SELASA, RABU, KAMIS, or JUMAT
Jam        : HH.MM - HH.MM (e.g., 07.00 - 07.50)
Ruangan    : Room code (e.g., IF-101, IF-AV)
Prodi      : Study program (IF, IUP, RKA, RPL)
Mata Kuliah: Course name
Semester   : 1-8
Kode Dosen : Lecturer initials (2 letters)
SKS        : Credit units
```

## Main Library

- [excelize](https://github.com/xuri/excelize) - Go library for reading & writing Excel files (.xlsx)
