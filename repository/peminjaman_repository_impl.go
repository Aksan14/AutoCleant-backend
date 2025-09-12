package repository

import (
	"context"
	"database/sql"
	"errors"
	"reset/model"
)

type peminjamanRepositoryImpl struct {
	DB *sql.DB
}

func NewPeminjamanRepositoryImpl(db *sql.DB) PeminjamanRepository {
	return &peminjamanRepositoryImpl{DB: db}
}

func (r *peminjamanRepositoryImpl) CreateTx(ctx context.Context, tx *sql.Tx, p model.Peminjaman) (model.Peminjaman, error) {
	res, err := tx.ExecContext(ctx, `
	   INSERT INTO peminjaman (inventaris_id, nama_peminjam, tgl_pinjam, rencana_kembali, jumlah, keterangan, status)
	   VALUES (?, ?, ?, ?, ?, ?, 'dipinjam')`,
		p.BarangID, p.NamaPeminjam, p.TglPinjam, p.RencanaKembali, p.Jumlah, p.Keterangan,
	)
	if err != nil {
		return p, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return p, err
	}
	p.ID = int(id)
	p.Status = "dipinjam"
	return p, nil
}

func (r *peminjamanRepositoryImpl) UpdateReturnTx(ctx context.Context, tx *sql.Tx, id int, tglKembali string, kondisi string, fotoBukti string, keteranganKembali string) error {
	res, err := tx.ExecContext(ctx, `
		UPDATE peminjaman
		SET tgl_kembali = ?, kondisi_setelah = ?, status = 'selesai', foto_bukti_kembali = ?, keterangan_kembali = ?
		WHERE id = ? AND status = 'dipinjam'`,
		tglKembali, kondisi, fotoBukti, keteranganKembali, id,
	)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return errors.New("data peminjaman tidak ditemukan atau sudah selesai")
	}
	return nil
}

func (r *peminjamanRepositoryImpl) GetAll(ctx context.Context) ([]model.Peminjaman, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT p.id, p.inventaris_id, b.nama_barang, p.nama_peminjam, 
		       p.tgl_pinjam, p.rencana_kembali, p.tgl_kembali, p.kondisi_setelah, p.status, p.keterangan,
		       p.foto_bukti_kembali, p.keterangan_kembali, p.jumlah
		FROM peminjaman p
		JOIN inventaris b ON b.id = p.inventaris_id
		ORDER BY p.tgl_pinjam DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Peminjaman
	for rows.Next() {
		var m model.Peminjaman
		var tglKembali sql.NullTime
		var kondisi sql.NullString
		var fotoBukti sql.NullString
		var keteranganKembali sql.NullString
		if err := rows.Scan(&m.ID, &m.BarangID, &m.BarangNama, &m.NamaPeminjam, &m.TglPinjam, &m.RencanaKembali, &tglKembali, &kondisi, &m.Status, &m.Keterangan, &fotoBukti, &keteranganKembali, &m.Jumlah); err != nil {
			return nil, err
		}
		if tglKembali.Valid {
			m.TglKembali = &tglKembali.Time
		}
		if kondisi.Valid {
			s := kondisi.String
			m.KondisiSetelah = &s
		}
		if fotoBukti.Valid {
			s := fotoBukti.String
			m.FotoBuktiKembali = &s
		}
		if keteranganKembali.Valid {
			s := keteranganKembali.String
			m.KeteranganKembali = &s
		}
		list = append(list, m)
	}
	return list, nil
}

func (r *peminjamanRepositoryImpl) GetByID(ctx context.Context, id int) (model.Peminjaman, error) {
	var m model.Peminjaman
	var tglKembali sql.NullTime
	var kondisi sql.NullString
	var fotoBukti sql.NullString
	var keteranganKembali sql.NullString
	err := r.DB.QueryRowContext(ctx, `
		SELECT id, inventaris_id, nama_peminjam, tgl_pinjam, rencana_kembali, tgl_kembali, kondisi_setelah, status, keterangan, foto_bukti_kembali, keterangan_kembali, jumlah
		FROM peminjaman WHERE id = ?`, id).
		Scan(&m.ID, &m.BarangID, &m.NamaPeminjam, &m.TglPinjam, &m.RencanaKembali, &tglKembali, &kondisi, &m.Status, &m.Keterangan, &fotoBukti, &keteranganKembali, &m.Jumlah)
	if err != nil {
		return model.Peminjaman{}, err
	}
	if tglKembali.Valid {
		m.TglKembali = &tglKembali.Time
	}
	if kondisi.Valid {
		s := kondisi.String
		m.KondisiSetelah = &s
	}
	if fotoBukti.Valid {
		s := fotoBukti.String
		m.FotoBuktiKembali = &s
	}
	if keteranganKembali.Valid {
		s := keteranganKembali.String
		m.KeteranganKembali = &s
	}
	return m, nil
}

func (r *peminjamanRepositoryImpl) DeleteByID(ctx context.Context, id int) error {
	res, err := r.DB.ExecContext(ctx, `
		DELETE FROM peminjaman WHERE id = ?`, id)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return errors.New("data peminjaman tidak ditemukan")
	}
	return nil
}
