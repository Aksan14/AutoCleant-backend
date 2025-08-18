package repository

import (
	"context"
	"database/sql"
	"reset/model"
)

type PeminjamanRepository interface {
	CreateTx(ctx context.Context, tx *sql.Tx, p model.Peminjaman) (model.Peminjaman, error)
	UpdateReturnTx(ctx context.Context, tx *sql.Tx, id int, tglKembali string, kondisi string) error
	GetAll(ctx context.Context) ([]model.Peminjaman, error)
	GetByID(ctx context.Context, id int) (model.Peminjaman, error)
}