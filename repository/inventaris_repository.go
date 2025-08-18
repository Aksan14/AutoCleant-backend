package repository

import (
    "database/sql"
    "errors"
    "reset/model"
    _ "github.com/go-sql-driver/mysql"
    "time"
)

type itemRepositoryMySQL struct {
    db *sql.DB
}

func NewItemRepositoryImpl(db *sql.DB) ItemRepository {
    return &itemRepositoryMySQL{db: db}
}

func (r *itemRepositoryMySQL) CreateInventaris(item *model.Inventaris) (int, error) {
    query := `INSERT INTO inventaris (nama_barang, kategori, jumlah, satuan, kondisi, foto) VALUES (?, ?, ?, ?, ?, ?)`
    res, err := r.db.Exec(query, item.NamaBarang, item.Kategori, item.Jumlah, item.Satuan, item.Kondisi, item.Foto)
    if err != nil {
        return 0, err
    }
    id, err := res.LastInsertId()
    return int(id), err
}

func (r *itemRepositoryMySQL) GetByIDInventaris(id int) (*model.Inventaris, error) {
    query := `SELECT id, nama_barang, kategori, jumlah, satuan, kondisi, foto FROM inventaris WHERE id = ?`
    row := r.db.QueryRow(query, id)
    it := &model.Inventaris{}
    if err := row.Scan(&it.ID, &it.NamaBarang, &it.Kategori, &it.Jumlah, &it.Satuan, &it.Kondisi, &it.Foto); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return it, nil
}

func (r *itemRepositoryMySQL) GetAllInventaris() ([]*model.Inventaris, error) {
    query := `SELECT id, nama_barang, kategori, jumlah, satuan, kondisi, foto FROM inventaris ORDER BY id DESC`
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []*model.Inventaris
    for rows.Next() {
        it := &model.Inventaris{}
        var created, updated time.Time
        if err := rows.Scan(&it.ID, &it.NamaBarang, &it.Kategori, &it.Jumlah, &it.Satuan, &it.Kondisi, &it.Foto); err != nil {
            return nil, err
        }
        it.CreatedAt = created
        it.UpdatedAt = updated
        items = append(items, it)
    }
    return items, nil
}

func (r *itemRepositoryMySQL) UpdateInventaris(item *model.Inventaris) error {
    query := `UPDATE inventaris SET nama_barang=?, kategori=?, jumlah=?, satuan=?, kondisi=?, foto=? WHERE id=?`
    _, err := r.db.Exec(query, item.NamaBarang, item.Kategori, item.Jumlah, item.Satuan, item.Kondisi, item.Foto, item.ID)
    return err
}

func (r *itemRepositoryMySQL) DeleteInventaris(id int, kondisi string) error {
    query := `UPDATE inventaris SET kondisi = ? WHERE id = ?`
    _, err := r.db.Exec(query, kondisi, id)
    return err
}

func (r *itemRepositoryMySQL) FindByNamabarang(nama string) (*model.Inventaris, error) {
	query := `SELECT id, nama_barang, kategori, jumlah, satuan, foto, kondisi
			FROM inventaris WHERE nama_barang = ? LIMIT 1`

	row := r.db.QueryRow(query, nama)

	var item model.Inventaris
	err := row.Scan(&item.ID, &item.NamaBarang, &item.Kategori, &item.Jumlah, &item.Satuan, &item.Foto, &item.Kondisi)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &item, nil
}

func (r *itemRepositoryMySQL) IsBarangSedangDipinjam(barangID int) (bool, error) {
    var count int
    query := `SELECT COUNT(*) FROM peminjaman WHERE inventaris_id = ?`
    err := r.db.QueryRow(query, barangID).Scan(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

func (r *itemRepositoryMySQL) SearchInventaris(keyword string) ([]model.Inventaris, error) {
    query := `
        SELECT id, nama_barang, kategori, jumlah, satuan, kondisi, foto
        FROM inventaris
        WHERE nama_barang LIKE CONCAT('%', ?, '%')
            OR kategori LIKE CONCAT('%', ?, '%')`
    
    rows, err := r.db.Query(query, keyword, keyword)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []model.Inventaris
    for rows.Next() {
        var item model.Inventaris
        if err := rows.Scan(&item.ID, &item.NamaBarang, &item.Kategori, &item.Jumlah, &item.Satuan, &item.Kondisi, &item.Foto); err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    return items, nil
}
