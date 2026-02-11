package main

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
