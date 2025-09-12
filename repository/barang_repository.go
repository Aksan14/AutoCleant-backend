package repository

import (
	"context"
	"database/sql"
	"reset/model"
)

type BarangRepository interface {
	GetAvailable(ctx context.Context) ([]model.Inventaris, error)
	GetByID(ctx context.Context, id int) (model.Inventaris, error)
	UpdateStatusTx(ctx context.Context, tx *sql.Tx, id int, status string) error
	UpdateJumlahTx(ctx context.Context, tx *sql.Tx, id int, jumlah int) error
}
