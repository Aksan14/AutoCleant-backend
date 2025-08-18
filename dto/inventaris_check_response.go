package dto

import "reset/model"

type StartReportResponse struct {
	ID         int    `json:"id"`
	KodeReport string `json:"kode_report"`
}

type CheckResponse struct {
	ID           int    `json:"id"`
	InventarisID int    `json:"inventaris_id"`
	Kondisi      string `json:"kondisi"`
	Keterangan   string `json:"keterangan"`
}

type ReportResponse struct {
	Report model.InventarisReport `json:"report"`
	Checks []model.InventarisCheck  `json:"checks"`
	Items  []model.Inventaris          `json:"items"`
}