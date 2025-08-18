package dto

type StartReportRequest struct {
	Petugas string `json:"petugas"`
}

type CheckRequest struct {
	InventarisID int    `json:"inventaris_id"`
	Kondisi      string `json:"kondisi"`
	Keterangan   string `json:"keterangan"`
}
