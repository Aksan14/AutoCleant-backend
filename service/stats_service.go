package service

import "context"

type ReportService interface {
	GetCountAllBarang(ctx context.Context) (int, error)
	GetCountBarangDipinjam(ctx context.Context) (int, error)
	GetCountBarangRusakBerat(ctx context.Context) (int, error)
}