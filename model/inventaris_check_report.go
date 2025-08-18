package model

import "time"

type InventarisReport struct {
	ID            int       `json:"id"`
	KodeReport    string    `json:"kode_report"`
	TanggalReport time.Time `json:"tanggal_report"`
	Petugas       string    `json:"petugas"`
	Status        string    `json:"status"`
}

type InventarisCheck struct {
	ID           int       `json:"id"`
	ReportID     int       `json:"report_id"`
	InventarisID int       `json:"inventaris_id"`
	Kondisi      string    `json:"kondisi"`
	Keterangan   string    `json:"keterangan"`
	TanggalCek   time.Time `json:"tanggal_cek"`
}