package service

import (
	"encoding/json"
	"fmt"

	"reset/dto"
	"reset/repository"
	"reset/model"
)

// ServiceImpl implements the Service interface
type ServiceImpl struct {
	repo repository.Repository
}

// NewService creates a new service instance
func NewService(repo repository.Repository) Service {
	return &ServiceImpl{repo: repo}
}

// -------------------- REPORT --------------------

func (s *ServiceImpl) StartReport(req dto.StartReportRequest) (dto.StartReportResponse, error) {
	report, err := s.repo.StartReport(req.Petugas)
	if err != nil {
		return dto.StartReportResponse{}, err
	}
	return dto.StartReportResponse{
		ID:         report.ID,
		KodeReport: report.KodeReport,
	}, nil
}

func (s *ServiceImpl) FinalizeReport(id int) error {
	return s.repo.FinalizeReport(id)
}

func (s *ServiceImpl) GetReports(status string) ([]model.InventarisReport, error) {
	return s.repo.GetReports(status)
}

func (s *ServiceImpl) GetReportDetail(id int) (dto.ReportResponse, error) {
	return s.repo.GetReportDetail(id)
}

// -------------------- CHECK --------------------

func (s *ServiceImpl) AddCheck(reportID int, req dto.CheckRequest) (dto.CheckResponse, error) {
	check, err := s.repo.AddCheck(reportID, req)
	if err != nil {
		return dto.CheckResponse{}, err
	}
	return dto.CheckResponse{
		ID:           check.ID,
		InventarisID: check.InventarisID,
		Kondisi:      check.Kondisi,
		Keterangan:   check.Keterangan,
	}, nil
}

func (s *ServiceImpl) UpdateCheck(id int, req dto.CheckRequest) error {
	return s.repo.UpdateCheck(id, req)
}

func (s *ServiceImpl) DeleteCheck(id int) error {
	return s.repo.DeleteCheck(id)
}

func (s *ServiceImpl) ExportPDF(id int) ([]byte, error) {
	detail, err := s.repo.GetReportDetail(id)
	if err != nil {
		return nil, err
	}
	if detail.Report.Status != "final" {
		return nil, fmt.Errorf("can only export finalized reports")
	}

	// Placeholder: nanti ganti dengan library PDF
	data, err := json.Marshal(detail)
	if err != nil {
		return nil, err
	}
	return data, nil
}
