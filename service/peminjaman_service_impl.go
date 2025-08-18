package service

import (
	"context"
	"database/sql"
	"errors"
	"reset/dto"
	"reset/model"
	"reset/repository"
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
			Satuan: b.Satuan, Kondisi: b.Kondisi, Foto: b.Foto,
		})
	}
	return out, nil
}

func (s *peminjamanServiceImpl) CreatePeminjaman(ctx context.Context, req dto.CreatePeminjamanRequest) (dto.PeminjamanResponse, error) {
	// Validasi barang tersedia
	barang, err := s.brgRepo.GetByID(ctx, req.BarangID)
	if err != nil {
		return dto.PeminjamanResponse{}, err
	}
	if barang.Status != "tersedia" {
		return dto.PeminjamanResponse{}, errors.New("barang tidak tersedia")
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
		Keterangan:    req.Keterangan,
	})
	if err != nil {
		return dto.PeminjamanResponse{}, err
	}

	if err = s.brgRepo.UpdateStatusTx(ctx, tx, req.BarangID, "dipinjam"); err != nil {
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
		Status:         "dipinjam",
		Keterangan:    p.Keterangan,
	}, nil
}

func (s *peminjamanServiceImpl) ReturnPeminjaman(ctx context.Context, id int, req dto.ReturnPeminjamanRequest) error {
	tKembali, err := time.Parse("2006-01-02", req.TglKembali)
	if err != nil {
		return errors.New("format tgl_kembali salah (YYYY-MM-DD)")
	}


	pjm, err := s.pjmRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if pjm.Status == "selesai" {
		return errors.New("peminjaman sudah selesai dan terkunci")
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

	// update peminjaman -> selesai
	if err = s.pjmRepo.UpdateReturnTx(ctx, tx, id, tKembali.Format("2006-01-02"), req.KondisiSetelah); err != nil {
		return err
	}
	// update status barang -> tersedia
	if err = s.brgRepo.UpdateStatusTx(ctx, tx, pjm.BarangID, "tersedia"); err != nil {
		return err
	}

	return tx.Commit()
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
			NamaPeminjam:   m.NamaPeminjam,
			TglPinjam:      m.TglPinjam.Format("2006-01-02"),
			RencanaKembali: m.RencanaKembali.Format("2006-01-02"),
			TglKembali:     tglKembali, KondisiSetelah: m.KondisiSetelah,
			Status: m.Status,
			Keterangan: m.Keterangan,
		})
	}
	return out, nil
}
