package repository

import "context"

type ReportRepository interface {
	CountAllBarang(ctx context.Context) (int, error)
	CountBarangDipinjam(ctx context.Context) (int, error)
	CountBarangRusakBerat(ctx context.Context) (int, error)
}