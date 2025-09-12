package service

import (
	"context"
	"database/sql"
	"errors"
	"reset/dto"
	"reset/model"
	"reset/repository"
	"strconv"
	"time"
)

type peminjamanServiceImpl struct {
	db      *sql.DB
	brgRepo repository.BarangRepository
	pjmRepo repository.PeminjamanRepository
}

func NewPeminjamanServiceImpl(db *sql.DB, br repository.BarangRepository, pr repository.PeminjamanRepository) PeminjamanService {
	return &peminjamanServiceImpl{db: db, brgRepo: br, pjmRepo: pr}
}

func (s *peminjamanServiceImpl) ListBarangTersedia(ctx context.Context) ([]dto.BarangSimpleResponse, error) {
	items, err := s.brgRepo.GetAvailable(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.BarangSimpleResponse, 0, len(items))
	for _, b := range items {
		out = append(out, dto.BarangSimpleResponse{
			ID: b.ID, NamaBarang: b.NamaBarang, Kategori: b.Kategori,
			Satuan: b.Satuan, Kondisi: b.Kondisi, Foto: b.Foto, Jumlah: b.Jumlah,
		})
	}
	return out, nil
}

func (s *peminjamanServiceImpl) CreatePeminjaman(ctx context.Context, req dto.CreatePeminjamanRequest) (dto.PeminjamanResponse, error) {
	barang, err := s.brgRepo.GetByID(ctx, req.BarangID)
	if err != nil {
		return dto.PeminjamanResponse{}, err
	}
	if barang.Status != "tersedia" {
		return dto.PeminjamanResponse{}, errors.New("barang tidak tersedia")
	}
	if req.Jumlah < 1 {
		return dto.PeminjamanResponse{}, errors.New("jumlah minimal 1")
	}
	if req.Jumlah > barang.Jumlah {
		return dto.PeminjamanResponse{}, errors.New("jumlah yang diminta melebihi stok tersedia")
	}

	tPinjam, err := time.Parse("2006-01-02", req.TglPinjam)
	if err != nil {
		return dto.PeminjamanResponse{}, errors.New("format tgl_pinjam salah (YYYY-MM-DD)")
	}
	tRencana, err := time.Parse("2006-01-02", req.RencanaKembali)
	if err != nil {
		return dto.PeminjamanResponse{}, errors.New("format rencana_kembali salah (YYYY-MM-DD)")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return dto.PeminjamanResponse{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// create peminjaman
	p, err := s.pjmRepo.CreateTx(ctx, tx, model.Peminjaman{
		BarangID:       req.BarangID,
		NamaPeminjam:   req.NamaPeminjam,
		TglPinjam:      tPinjam,
		RencanaKembali: tRencana,
		Jumlah:         req.Jumlah,
		Keterangan:     req.Keterangan,
	})
	if err != nil {
		return dto.PeminjamanResponse{}, err
	}

	// Kurangi stok barang
	newJumlah := barang.Jumlah - req.Jumlah
	err = s.brgRepo.UpdateJumlahTx(ctx, tx, req.BarangID, newJumlah)
	if err != nil {
		return dto.PeminjamanResponse{}, err
	}

	// Update status otomatis
	status := "tersedia"
	if newJumlah == 0 {
		status = "dipinjam"
	}
	if err = s.brgRepo.UpdateStatusTx(ctx, tx, req.BarangID, status); err != nil {
		return dto.PeminjamanResponse{}, err
	}

	if err = tx.Commit(); err != nil {
		return dto.PeminjamanResponse{}, err
	}

	return dto.PeminjamanResponse{
		ID: p.ID, BarangID: p.BarangID, BarangNama: barang.NamaBarang,
		NamaPeminjam:   p.NamaPeminjam,
		TglPinjam:      tPinjam.Format("2006-01-02"),
		RencanaKembali: tRencana.Format("2006-01-02"),
		Jumlah:         p.Jumlah,
		Status:         p.Status,
		Keterangan:     p.Keterangan,
	}, nil
}

func (s *peminjamanServiceImpl) ReturnPeminjaman(ctx context.Context, id int, req dto.ReturnPeminjamanRequest) error {
	// Validasi input required
	if req.FotoBuktiKembali == "" {
		return errors.New("foto bukti pengembalian wajib diisi")
	}
	if req.KeteranganKembali == "" {
		return errors.New("keterangan pengembalian wajib diisi")
	}

	tKembali, err := time.Parse("2006-01-02", req.TglKembali)
	if err != nil {
		return errors.New("format tgl_kembali salah (YYYY-MM-DD)")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1. Ambil data peminjaman dengan jumlah yang dipinjam
	pjm, err := s.pjmRepo.GetByID(ctx, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 2. Validasi status peminjaman
	if pjm.Status == "selesai" {
		tx.Rollback()
		return errors.New("peminjaman sudah selesai dan terkunci")
	}

	// 3. Validasi jumlah yang dipinjam harus valid
	if pjm.Jumlah <= 0 {
		tx.Rollback()
		return errors.New("data peminjaman tidak valid, jumlah kosong")
	}

	// 4. Ambil data barang dari inventaris
	barang, err := s.brgRepo.GetByID(ctx, pjm.BarangID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 5. Kembalikan stok barang ke inventaris
	// Tambahkan jumlah yang dipinjam kembali ke stok
	newJumlah := barang.Jumlah + pjm.Jumlah
	if err = s.brgRepo.UpdateJumlahTx(ctx, tx, pjm.BarangID, newJumlah); err != nil {
		tx.Rollback()
		return err
	}

	// 6. Update status barang otomatis
	status := "tersedia"
	if newJumlah == 0 {
		status = "dipinjam"
	}
	if err = s.brgRepo.UpdateStatusTx(ctx, tx, pjm.BarangID, status); err != nil {
		tx.Rollback()
		return err
	}

	// 7. Update data peminjaman (status, tanggal kembali, kondisi, foto, keterangan)
	// PENTING: Field Jumlah TIDAK BERUBAH - tetap untuk history
	if err = s.pjmRepo.UpdateReturnTx(ctx, tx, id, tKembali.Format("2006-01-02"), req.KondisiSetelah, req.FotoBuktiKembali, req.KeteranganKembali); err != nil {
		tx.Rollback()
		return err
	}

	// 8. Commit transaksi
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *peminjamanServiceImpl) ListPeminjaman(ctx context.Context) ([]dto.PeminjamanResponse, error) {
	list, err := s.pjmRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.PeminjamanResponse, 0, len(list))
	for _, m := range list {
		var tglKembali *string
		if m.TglKembali != nil {
			s := m.TglKembali.Format("2006-01-02")
			tglKembali = &s
		}
		out = append(out, dto.PeminjamanResponse{
			ID: m.ID, BarangID: m.BarangID, BarangNama: m.BarangNama,
			NamaPeminjam:      m.NamaPeminjam,
			TglPinjam:         m.TglPinjam.Format("2006-01-02"),
			RencanaKembali:    m.RencanaKembali.Format("2006-01-02"),
			TglKembali:        tglKembali,
			Jumlah:            m.Jumlah, 
			KondisiSetelah:    m.KondisiSetelah,
			Status:            m.Status,
			Keterangan:        m.Keterangan,
			FotoBuktiKembali:  m.FotoBuktiKembali,
			KeteranganKembali: m.KeteranganKembali,
		})
	}
	return out, nil
}

func (s *peminjamanServiceImpl) DeleteByID(ctx context.Context, id int) error {
	peminjaman, err := s.pjmRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("data tidak ditemukan dengan id " + strconv.Itoa(id))
	}

	if peminjaman.Status != "selesai" {
		return errors.New("tidak dapat menghapus peminjaman yang masih aktif")
	}

	return s.pjmRepo.DeleteByID(ctx, id)
}
