package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/xuri/excelize/v2"
)

// JadwalKuliah menyimpan data jadwal per mata kuliah
type JadwalKuliah struct {
	Hari       string `json:"hari"`
	Jam        string `json:"jam"`
	Ruangan    string `json:"ruangan"`
	Prodi      string `json:"prodi"`
	MataKuliah string `json:"mata_kuliah"`
	Semester   string `json:"semester"`
	KodeDosen  string `json:"kode_dosen"`
	SKS        string `json:"sks"`
	RawData    string `json:"raw_data"`
}

func main() {
	// Ganti dengan path file Excel kamu
	filePath := "data.xlsx"

	// Buka file Excel
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatalf("Gagal membuka file: %v", err)
	}
	defer f.Close()

	// Ambil nama sheet "Jadwal Kuliah"
	sheetName := "Jadwal Kuliah"
	fmt.Printf("Membaca sheet: %s\n\n", sheetName)

	// Ambil semua baris dari sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatalf("Gagal membaca rows: %v", err)
	}

	if len(rows) < 2 {
		log.Fatal("File Excel kosong atau tidak ada data")
	}

	// Baris 1 (index 0): Header ruangan (IF-101, IF-102, dll)
	headerRuangan := rows[0]

	// Mapping jadwal
	var semuaJadwal []JadwalKuliah
	jadwalBySemester := make(map[string][]JadwalKuliah)
	jadwalByProdi := make(map[string][]JadwalKuliah)

	// Mapping hari
	hariList := []string{"SENIN", "SELASA", "RABU", "KAMIS", "JUMAT"}
	currentHari := ""

	// Mulai dari baris ke-2 (index 1)
	for rowIdx := 1; rowIdx < len(rows); rowIdx++ {
		row := rows[rowIdx]
		if len(row) == 0 {
			continue
		}

		// Kolom A: Hari (jika ada, update currentHari)
		if len(row) > 0 {
			cellA := strings.ToUpper(strings.TrimSpace(row[0]))
			for _, h := range hariList {
				if cellA == h {
					currentHari = h
					break
				}
			}
		}

		// Kolom B: Jam
		jamBaris1 := ""
		if len(row) > 1 {
			jamBaris1 = strings.TrimSpace(row[1])
		}

		if jamBaris1 == "" || currentHari == "" {
			continue
		}

		// Kolom G ke atas (index 6+): Data mata kuliah per ruangan
		for colIdx := 6; colIdx < len(row); colIdx++ {
			cellValue := strings.TrimSpace(row[colIdx])
			if cellValue == "" {
				continue
			}

			// Cek apakah ini baris nama matkul
			// Jika sesuai format PRODI_NamaMatkul atau hanya nama matkul (akan jadi S2)

			// Ini baris nama matkul, ambil baris berikutnya untuk info sem/dosen/sks
			semInfo := ""
			jamBaris2 := ""
			if rowIdx+1 < len(rows) {
				nextRow := rows[rowIdx+1]
				if colIdx < len(nextRow) {
					semInfo = strings.TrimSpace(nextRow[colIdx])
				}
				// Ambil jam dari baris berikutnya
				if len(nextRow) > 1 {
					jamBaris2 = strings.TrimSpace(nextRow[1])
				}
			}

			// Extract jam mulai dan jam selesai, lalu gabungkan
			jamMulai := extractJamMulai(jamBaris1)
			jamSelesai := extractJamSelesai(jamBaris2)
			jam := strings.TrimSpace(jamMulai)
			if jamSelesai != "" {
				jam = fmt.Sprintf("%s - %s", jamMulai, jamSelesai)
			}

			// Ambil nama ruangan dari header dan bersihkan
			ruangan := ""
			if colIdx < len(headerRuangan) {
				ruangan = cleanRuangan(strings.TrimSpace(headerRuangan[colIdx]))
			}

			// Parse data mata kuliah
			rawData := cellValue + "\n" + semInfo
			jadwal := parseJadwal(currentHari, jam, ruangan, cellValue, semInfo, rawData)

			if jadwal.MataKuliah != "" {
				semuaJadwal = append(semuaJadwal, jadwal)

				// Group by semester
				if jadwal.Semester != "" {
					jadwalBySemester[jadwal.Semester] = append(jadwalBySemester[jadwal.Semester], jadwal)
				}

				// Group by prodi
				if jadwal.Prodi != "" {
					jadwalByProdi[jadwal.Prodi] = append(jadwalByProdi[jadwal.Prodi], jadwal)
				}
			}
		}
	}

	// Output hasil
	fmt.Println("=== SEMUA JADWAL ===")
	fmt.Printf("Total: %d jadwal ditemukan\n\n", len(semuaJadwal))

	fmt.Println("=== JADWAL PER SEMESTER ===")
	for sem, jadwals := range jadwalBySemester {
		fmt.Printf("  Semester %s: %d mata kuliah\n", sem, len(jadwals))
	}

	fmt.Println("\n=== JADWAL PER PRODI ===")
	for prodi, jadwals := range jadwalByProdi {
		fmt.Printf("  %s: %d mata kuliah\n", prodi, len(jadwals))
	}

	// Export ke JSON file
	result := map[string]interface{}{
		"jadwal_by_semester": jadwalBySemester,
		"jadwal_by_prodi":    jadwalByProdi,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Gagal convert ke JSON: %v", err)
	}

	// Simpan ke file
	outputFile := "jadwal.json"
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		log.Fatalf("Gagal menulis file: %v", err)
	}

	fmt.Printf("\n✓ Berhasil export ke %s\n", outputFile)

	// Export ke file Excel
	fmt.Println("\n=== EXPORT KE EXCEL ===")
	err = exportToExcel(jadwalByProdi)
	if err != nil {
		log.Fatalf("Gagal export ke Excel: %v", err)
	}
}

// exportToExcel menulis data ke file Excel lokal
func exportToExcel(jadwalByProdi map[string][]JadwalKuliah) error {
	outputFile := "jadwal.xlsx"
	f := excelize.NewFile()

	defaultSheet := f.GetSheetName(0)

	// Export jadwal per prodi dan semester
	fmt.Println("\nExport jadwal per prodi dan semester...")
	jadwalByProdiSemester := make(map[string]map[string][]JadwalKuliah)
	for prodi, jadwals := range jadwalByProdi {
		for _, j := range jadwals {
			if j.Semester == "" {
				continue
			}
			if _, ok := jadwalByProdiSemester[prodi]; !ok {
				jadwalByProdiSemester[prodi] = make(map[string][]JadwalKuliah)
			}
			jadwalByProdiSemester[prodi][j.Semester] = append(jadwalByProdiSemester[prodi][j.Semester], j)
		}
	}

	var prodis []string
	for prodi := range jadwalByProdiSemester {
		prodis = append(prodis, prodi)
	}
	sort.Strings(prodis)

	for _, prodi := range prodis {
		var semesters []string
		for sem := range jadwalByProdiSemester[prodi] {
			semesters = append(semesters, sem)
		}
		sort.Strings(semesters)

		for _, sem := range semesters {
			jadwals := jadwalByProdiSemester[prodi][sem]
			sheetTitle := fmt.Sprintf("%s_Sem_%s", prodi, sem)
			if _, err := f.NewSheet(sheetTitle); err != nil {
				return fmt.Errorf("gagal membuat sheet %s: %v", sheetTitle, err)
			}
			if err := writeSheetData(f, sheetTitle, jadwals); err != nil {
				return fmt.Errorf("gagal menulis sheet %s: %v", sheetTitle, err)
			}
			fmt.Printf("  ✓ Sheet '%s' (%d data)\n", sheetTitle, len(jadwals))
		}
	}

	// Hapus default sheet kalau masih ada dan tidak dipakai
	if defaultSheet != "" {
		_ = f.DeleteSheet(defaultSheet)
	}

	if err := f.SaveAs(outputFile); err != nil {
		return fmt.Errorf("gagal menyimpan file excel: %v", err)
	}

	fmt.Printf("\n✓ Berhasil export ke %s\n", outputFile)
	return nil
}

// writeSheetData menulis header dan data jadwal ke sheet
func writeSheetData(f *excelize.File, sheetTitle string, jadwals []JadwalKuliah) error {
	// Header
	header := []string{
		"Hari", "Jam", "Ruangan", "Prodi", "Mata Kuliah", "Semester", "Kode Dosen", "SKS",
	}
	for i, h := range header {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(sheetTitle, cell, h); err != nil {
			return err
		}
	}

	// Data rows
	for rowIdx, j := range jadwals {
		r := rowIdx + 2
		values := []string{
			j.Hari,
			j.Jam,
			j.Ruangan,
			j.Prodi,
			j.MataKuliah,
			j.Semester,
			j.KodeDosen,
			j.SKS,
		}
		for colIdx, v := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, r)
			if err := f.SetCellValue(sheetTitle, cell, v); err != nil {
				return err
			}
		}
	}

	return nil
}

// cleanRuangan membersihkan nama ruangan dari suffix seperti "a&b"
// "IF-105a&b (kapasitas 100)" -> "IF-105 (kapasitas 100)"
func cleanRuangan(ruangan string) string {
	// Regex untuk menghapus suffix seperti "a&b", "a&b&c", dll sebelum spasi atau kurung
	re := regexp.MustCompile(`([A-Z]+-\d+)[a-z&]+\s*(\(.*\))?`)
	if re.MatchString(ruangan) {
		return re.ReplaceAllString(ruangan, "$1 $2")
	}
	return ruangan
}

// extractJamMulai mengambil jam mulai dari format "07.00 - 07.50"
func extractJamMulai(jam string) string {
	jam = strings.ReplaceAll(jam, " ", "")
	if parts := strings.Split(jam, "-"); len(parts) >= 1 {
		return strings.TrimSpace(parts[0])
	}
	return jam
}

// extractJamSelesai mengambil jam selesai dari format "08.00 - 08.50"
func extractJamSelesai(jam string) string {
	jam = strings.ReplaceAll(jam, " ", "")
	if parts := strings.Split(jam, "-"); len(parts) == 2 {
		return strings.TrimSpace(parts[1])
	}
	return jam
}

// parseJadwal mengekstrak informasi dari cell jadwal
func parseJadwal(hari, jam, ruangan, matkulCell, semInfoCell, rawData string) JadwalKuliah {
	jadwal := JadwalKuliah{
		Hari:    hari,
		Jam:     jam,
		Ruangan: ruangan,
		RawData: rawData,
	}

	// Parse baris 1: PRODI_NamaMataKuliah
	prodiRegex := regexp.MustCompile(`^(IF|IUP|RKA|RPL|S3)_(.+)$`)
	prodiMatch := prodiRegex.FindStringSubmatch(matkulCell)
	if len(prodiMatch) == 3 {
		jadwal.Prodi = prodiMatch[1]
		jadwal.MataKuliah = strings.TrimSpace(prodiMatch[2])
	} else {
		// Jika tidak sesuai format standar, anggap sebagai mata kuliah S2
		jadwal.Prodi = "S2"
		jadwal.MataKuliah = strings.TrimSpace(matkulCell)
	}

	// Parse baris 2: Sem X / KODE_DOSEN / Y SKS
	// Format: "Sem 7 / SL / 3 SKS" atau "Sem 4 / SR, MA / 3 SKS"
	if semInfoCell != "" {
		parts := strings.Split(semInfoCell, "/")
		if len(parts) >= 3 {
			// Part 1: Semester
			semPart := strings.TrimSpace(parts[0])
			semRegex := regexp.MustCompile(`(?i)sem\s*(\d+)`)
			semMatch := semRegex.FindStringSubmatch(semPart)
			if len(semMatch) > 1 {
				jadwal.Semester = semMatch[1]
			}

			// Part 2: Kode Dosen (2 huruf, bisa multiple dipisah koma)
			jadwal.KodeDosen = strings.TrimSpace(parts[1])

			// Part 3: SKS
			sksPart := strings.TrimSpace(parts[2])
			sksRegex := regexp.MustCompile(`(\d+)\s*SKS`)
			sksMatch := sksRegex.FindStringSubmatch(sksPart)
			if len(sksMatch) > 1 {
				jadwal.SKS = sksMatch[1]
			}
		}
	}

	return jadwal
}
