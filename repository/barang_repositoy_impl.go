package repository

import (
	"context"
	"database/sql"
	"errors"
	"reset/model"
)

type barangRepositoryImpl struct {
	DB *sql.DB
}

func NewBarangRepositoryImpl(db *sql.DB) BarangRepository {
	return &barangRepositoryImpl{DB: db}
}

func (r *barangRepositoryImpl) GetAvailable(ctx context.Context) ([]model.Inventaris, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, nama_barang, kategori, jumlah, satuan, kondisi, foto, status
		 FROM inventaris 
		 WHERE status = 'tersedia' 
		   AND kondisi <> 'Dimusnahkan'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Inventaris
	for rows.Next() {
		var b model.Inventaris
		if err := rows.Scan(&b.ID, &b.NamaBarang, &b.Kategori, &b.Jumlah, &b.Satuan, &b.Kondisi, &b.Foto, &b.Status); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, nil
}


func (r *barangRepositoryImpl) GetByID(ctx context.Context, id int) (model.Inventaris, error) {
	var b model.Inventaris
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, nama_barang, kategori, jumlah, satuan, kondisi, foto, status
		 FROM inventaris WHERE id = ?`, id).
		Scan(&b.ID, &b.NamaBarang, &b.Kategori, &b.Jumlah, &b.Satuan, &b.Kondisi, &b.Foto, &b.Status)
	if err != nil {
		return model.Inventaris{}, err
	}
	return b, nil
}

func (r *barangRepositoryImpl) UpdateStatusTx(ctx context.Context, tx *sql.Tx, id int, status string) error {
	res, err := tx.ExecContext(ctx, `UPDATE inventaris SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return errors.New("barang tidak ditemukan")
	}
	return nil
}
