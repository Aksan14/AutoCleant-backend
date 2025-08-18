package service

import (
	"context"
	"reset/repository"
)

type reportServiceImpl struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportServiceImpl{repo: repo}
}

func (s *reportServiceImpl) GetCountAllBarang(ctx context.Context) (int, error) {
	return s.repo.CountAllBarang(ctx)
}

func (s *reportServiceImpl) GetCountBarangDipinjam(ctx context.Context) (int, error) {
	return s.repo.CountBarangDipinjam(ctx)
}

func (s *reportServiceImpl) GetCountBarangRusakBerat(ctx context.Context) (int, error) {
	return s.repo.CountBarangRusakBerat(ctx)
}