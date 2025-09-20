package repository

import (
	"database/sql"
	"fmt"
	"reset/model"
)

type BarangRepository interface {
	Create(db *sql.DB, barang *model.Barang) error
	GetByID(db *sql.DB, id int) (*model.Barang, error)
	GetAll(db *sql.DB, search string, limit, offset int) ([]model.Barang, error)
	GetTotalCount(db *sql.DB, search string) (int, error)
	Update(db *sql.DB, barang *model.Barang) error
	Delete(db *sql.DB, id int) error
}

type barangRepositoryImpl struct{}

func NewBarangRepository() BarangRepository {
	return &barangRepositoryImpl{}
}

func (r *barangRepositoryImpl) Create(db *sql.DB, barang *model.Barang) error {
	query := `INSERT INTO barang (nama_barang, deskripsi, gambar, harga, link_shopee, link_tiktokshop) 
              VALUES (?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(query,
		barang.NamaBarang,
		barang.Deskripsi,
		barang.Gambar,
		barang.Harga,
		barang.LinkShopee,
		barang.LinkTiktokshop,
	)
	if err != nil {
		return fmt.Errorf("gagal create barang: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("gagal get last insert id: %v", err)
	}

	barang.ID = int(id)
	return nil
}

func (r *barangRepositoryImpl) GetByID(db *sql.DB, id int) (*model.Barang, error) {
	query := `SELECT id, nama_barang, deskripsi, gambar, harga, link_shopee, link_tiktokshop, 
              created_at, updated_at 
              FROM barang WHERE id = ?`

	row := db.QueryRow(query, id)

	var barang model.Barang
	err := row.Scan(
		&barang.ID,
		&barang.NamaBarang,
		&barang.Deskripsi,
		&barang.Gambar,
		&barang.Harga,
		&barang.LinkShopee,
		&barang.LinkTiktokshop,
		&barang.CreatedAt,
		&barang.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("barang dengan id %d tidak ditemukan", id)
		}
		return nil, fmt.Errorf("gagal get barang: %v", err)
	}

	return &barang, nil
}

func (r *barangRepositoryImpl) GetAll(db *sql.DB, search string, limit, offset int) ([]model.Barang, error) {
	query := `SELECT id, nama_barang, deskripsi, gambar, harga, link_shopee, link_tiktokshop, 
              created_at, updated_at FROM barang`

	args := []interface{}{}

	if search != "" {
		query += " WHERE nama_barang LIKE ? OR deskripsi LIKE ?"
		searchParam := "%" + search + "%"
		args = append(args, searchParam, searchParam)
	}

	query += " ORDER BY created_at DESC"

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)

		if offset > 0 {
			query += " OFFSET ?"
			args = append(args, offset)
		}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("gagal query barang: %v", err)
	}
	defer rows.Close()

	var barangList []model.Barang
	for rows.Next() {
		var barang model.Barang
		err := rows.Scan(
			&barang.ID,
			&barang.NamaBarang,
			&barang.Deskripsi,
			&barang.Gambar,
			&barang.Harga,
			&barang.LinkShopee,
			&barang.LinkTiktokshop,
			&barang.CreatedAt,
			&barang.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan barang: %v", err)
		}
		barangList = append(barangList, barang)
	}

	return barangList, nil
}

func (r *barangRepositoryImpl) GetTotalCount(db *sql.DB, search string) (int, error) {
	query := "SELECT COUNT(*) FROM barang"
	args := []interface{}{}

	if search != "" {
		query += " WHERE nama_barang LIKE ? OR deskripsi LIKE ?"
		searchParam := "%" + search + "%"
		args = append(args, searchParam, searchParam)
	}

	var count int
	err := db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("gagal count barang: %v", err)
	}

	return count, nil
}

func (r *barangRepositoryImpl) Update(db *sql.DB, barang *model.Barang) error {
	query := `UPDATE barang SET 
              nama_barang = ?, 
              deskripsi = ?, 
              gambar = ?,
              harga = ?,
              link_shopee = ?, 
              link_tiktokshop = ?
              WHERE id = ?`

	result, err := db.Exec(query,
		barang.NamaBarang,
		barang.Deskripsi,
		barang.Gambar,
		barang.Harga,
		barang.LinkShopee,
		barang.LinkTiktokshop,
		barang.ID,
	)
	if err != nil {
		return fmt.Errorf("gagal update barang: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("barang dengan id %d tidak ditemukan", barang.ID)
	}

	return nil
}

func (r *barangRepositoryImpl) Delete(db *sql.DB, id int) error {
	query := "DELETE FROM barang WHERE id = ?"

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("gagal delete barang: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("barang dengan id %d tidak ditemukan", id)
	}

	return nil
}
