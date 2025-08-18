package service

import (
	"reset/dto"
	"reset/model"
)

type Service interface {
	StartReport(req dto.StartReportRequest) (dto.StartReportResponse, error)
	AddCheck(reportID int, req dto.CheckRequest) (dto.CheckResponse, error)
	UpdateCheck(id int, req dto.CheckRequest) error
	DeleteCheck(id int) error
	FinalizeReport(id int) error
	GetReports(status string) ([]model.InventarisReport, error)
	GetReportDetail(id int) (dto.ReportResponse, error)
	ExportPDF(id int) ([]byte, error)
}