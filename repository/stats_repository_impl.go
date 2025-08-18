package repository

import (
	"context"
	"database/sql"
)

type reportRepositoryImpl struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) ReportRepository {
	return &reportRepositoryImpl{db: db}
}

func (r *reportRepositoryImpl) CountAllBarang(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM inventaris").Scan(&count)
	return count, err
}

func (r *reportRepositoryImpl) CountBarangDipinjam(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM inventaris WHERE status = 'dipinjam'").Scan(&count)
	return count, err
}

func (r *reportRepositoryImpl) CountBarangRusakBerat(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM inventaris WHERE kondisi = 'Dimusnahkan'").Scan(&count)
	return count, err
}
