package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"reset/dto"
	"reset/model"

	"github.com/google/uuid"
)

type RepositoryImpl struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &RepositoryImpl{db: db}
}


func (r *RepositoryImpl) StartReport(petugas string) (model.InventarisReport, error) {
	kodeReport := "RPT-" + uuid.New().String()[:8]
	tanggal := time.Now()

	query := "INSERT INTO inventaris_report (kode_report, petugas, status, tanggal_report) VALUES (?, ?, 'draft', ?)"
	result, err := r.db.Exec(query, kodeReport, petugas, tanggal)
	if err != nil {
		return model.InventarisReport{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.InventarisReport{}, err
	}

	return model.InventarisReport{
		ID:            int(id),
		KodeReport:    kodeReport,
		Petugas:       petugas,
		Status:        "draft",
		TanggalReport: tanggal,
	}, nil
}

func (r *RepositoryImpl) FinalizeReport(id int) error {
	query := "UPDATE inventaris_report SET status = 'final' WHERE id = ? AND status = 'draft'"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("report not found or already finalized")
	}
	return nil
}

func (r *RepositoryImpl) GetReports(status string) ([]model.InventarisReport, error) {
	var rows *sql.Rows
	var err error

	if status == "" {
		rows, err = r.db.Query("SELECT id, kode_report, tanggal_report, petugas, status FROM inventaris_report ORDER BY id DESC;")
	} else {
		rows, err = r.db.Query("SELECT id, kode_report, tanggal_report, petugas, status FROM inventaris_report WHERE status = ? ORDER BY id DESC;", status)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []model.InventarisReport
	for rows.Next() {
		var r model.InventarisReport
		if err := rows.Scan(&r.ID, &r.KodeReport, &r.TanggalReport, &r.Petugas, &r.Status); err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}
	return reports, nil
}


func (r *RepositoryImpl) AddCheck(reportID int, check dto.CheckRequest) (model.InventarisCheck, error) {
    var status string
    err := r.db.QueryRow("SELECT status FROM inventaris_report WHERE id = ?", reportID).Scan(&status)
    if err == sql.ErrNoRows {
        return model.InventarisCheck{}, fmt.Errorf("report with id %d not found", reportID)
    }
    if err != nil {
        return model.InventarisCheck{}, err
    }

    if status != "draft" {
        return model.InventarisCheck{}, fmt.Errorf("cannot add check to finalized report")
    }

    tanggal := time.Now()
    query := "INSERT INTO inventaris_check (report_id, inventaris_id, kondisi, keterangan, tanggal_cek) VALUES (?, ?, ?, ?, ?)"
    result, err := r.db.Exec(query, reportID, check.InventarisID, check.Kondisi, check.Keterangan, tanggal)
    if err != nil {
        return model.InventarisCheck{}, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return model.InventarisCheck{}, err
    }

    return model.InventarisCheck{
        ID:           int(id),
        ReportID:     reportID,
        InventarisID: check.InventarisID,
        Kondisi:      check.Kondisi,
        Keterangan:   check.Keterangan,
        TanggalCek:   tanggal,
    }, nil
}


func (r *RepositoryImpl) UpdateCheck(id int, check dto.CheckRequest) error {
	var status string
	err := r.db.QueryRow("SELECT status FROM inventaris_report WHERE id = (SELECT report_id FROM inventaris_check WHERE id = ?)", id).Scan(&status)
	if err != nil {
		return err
	}
	if status != "draft" {
		return fmt.Errorf("cannot update check in finalized report")
	}

	query := "UPDATE inventaris_check SET inventaris_id = ?, kondisi = ?, keterangan = ? WHERE id = ?"
	_, err = r.db.Exec(query, check.InventarisID, check.Kondisi, check.Keterangan, id)
	return err
}

func (r *RepositoryImpl) DeleteCheck(id int) error {
	var status string
	err := r.db.QueryRow("SELECT status FROM inventaris_report WHERE id = (SELECT report_id FROM inventaris_check WHERE id = ?)", id).Scan(&status)
	if err != nil {
		return err
	}
	if status != "draft" {
		return fmt.Errorf("cannot delete check from finalized report")
	}

	query := "DELETE FROM inventaris_check WHERE id = ?"
	_, err = r.db.Exec(query, id)
	return err
}

func (r *RepositoryImpl) GetReportDetail(id int) (dto.ReportResponse, error) {
	var report model.InventarisReport
	err := r.db.QueryRow("SELECT id, kode_report, tanggal_report, petugas, status FROM inventaris_report WHERE id = ?", id).
		Scan(&report.ID, &report.KodeReport, &report.TanggalReport, &report.Petugas, &report.Status)
	if err != nil {
		return dto.ReportResponse{}, err
	}

	// Ambil semua checks
	rows, err := r.db.Query("SELECT id, report_id, inventaris_id, kondisi, keterangan, tanggal_cek FROM inventaris_check WHERE report_id = ?", id)
	if err != nil {
		return dto.ReportResponse{}, err
	}
	defer rows.Close()

	var checks []model.InventarisCheck
	var inventarisIDs []int
	inventarisMap := map[int]bool{}
	for rows.Next() {
		var c model.InventarisCheck
		if err := rows.Scan(&c.ID, &c.ReportID, &c.InventarisID, &c.Kondisi, &c.Keterangan, &c.TanggalCek); err != nil {
			return dto.ReportResponse{}, err
		}
		checks = append(checks, c)
		if !inventarisMap[c.InventarisID] {
			inventarisIDs = append(inventarisIDs, c.InventarisID)
			inventarisMap[c.InventarisID] = true
		}
	}

	items := []model.Inventaris{}
	if len(inventarisIDs) > 0 {
		placeholders := strings.TrimRight(strings.Repeat("?,", len(inventarisIDs)), ",")
		args := make([]interface{}, len(inventarisIDs))
		for i, id := range inventarisIDs {
			args[i] = id
		}
		query := fmt.Sprintf("SELECT id, nama_barang, kategori, jumlah, satuan, kondisi, foto, status FROM inventaris WHERE id IN (%s)", placeholders)
		rows2, err := r.db.Query(query, args...)
		if err != nil {
			return dto.ReportResponse{}, err
		}
		defer rows2.Close()

		for rows2.Next() {
			var item model.Inventaris
			if err := rows2.Scan(&item.ID, &item.NamaBarang, &item.Kategori, &item.Jumlah, &item.Satuan, &item.Kondisi, &item.Foto, &item.Status); err != nil {
				return dto.ReportResponse{}, err
			}
			items = append(items, item)
		}
	}

	return dto.ReportResponse{
		Report: report,
		Checks: checks,
		Items:  items,
	}, nil
}
