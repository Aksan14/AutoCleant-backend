package repository

import (
	"reset/dto"
	"reset/model"
)

// Repository defines the interface for database operations
type Repository interface {
	StartReport(petugas string) (model.InventarisReport, error)
	AddCheck(reportID int, check dto.CheckRequest) (model.InventarisCheck, error)
	UpdateCheck(id int, check dto.CheckRequest) error
	DeleteCheck(id int) error
	FinalizeReport(id int) error
	GetReports(status string) ([]model.InventarisReport, error)
	GetReportDetail(id int) (dto.ReportResponse, error)
}