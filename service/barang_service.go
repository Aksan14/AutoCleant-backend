package service

import (
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"reset/dto"
	"reset/model"
	"reset/repository"
	"strconv"
	"strings"
	"time"
)

type BarangService interface {
	CreateBarang(req dto.BarangFormRequest) (*dto.BarangResponse, error)
	GetBarangByID(id int) (*dto.BarangResponse, error)
	GetAllBarang(search string, page, limit int) (*dto.BarangListResponse, error)
	UpdateBarang(id int, req dto.BarangFormRequest) error
	DeleteBarang(id int) error
}

type barangServiceImpl struct {
	db   *sql.DB
	repo repository.BarangRepository
}

func NewBarangService(db *sql.DB, repo repository.BarangRepository) BarangService {
	return &barangServiceImpl{
		db:   db,
		repo: repo,
	}
}

func (s *barangServiceImpl) CreateBarang(req dto.BarangFormRequest) (*dto.BarangResponse, error) {
	// Validasi dan konversi harga
	harga, err := strconv.ParseFloat(req.Harga, 64)
	if err != nil {
		return nil, fmt.Errorf("format harga tidak valid")
	}

	if harga < 0 {
		return nil, fmt.Errorf("harga tidak boleh negatif")
	}

	var gambarPath string

	// Handle file upload if exists
	if req.Gambar != nil {
		filename, err := s.saveUploadedFile(req.Gambar)
		if err != nil {
			return nil, fmt.Errorf("gagal upload gambar: %v", err)
		}
		gambarPath = filename
	}

	barang := &model.Barang{
		NamaBarang:     req.NamaBarang,
		Deskripsi:      req.Deskripsi,
		Gambar:         gambarPath,
		Harga:          harga,
		LinkShopee:     req.LinkShopee,
		LinkTiktokshop: req.LinkTiktokshop,
	}

	if err := s.repo.Create(s.db, barang); err != nil {
		// Delete uploaded file if database save fails
		if gambarPath != "" {
			os.Remove(filepath.Join("uploads", "barang", gambarPath))
		}
		return nil, err
	}

	// Get created barang to return with timestamps
	createdBarang, err := s.repo.GetByID(s.db, barang.ID)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(createdBarang), nil
}

func (s *barangServiceImpl) GetBarangByID(id int) (*dto.BarangResponse, error) {
	barang, err := s.repo.GetByID(s.db, id)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(barang), nil
}

func (s *barangServiceImpl) GetAllBarang(search string, page, limit int) (*dto.BarangListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	barangList, err := s.repo.GetAll(s.db, search, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.GetTotalCount(s.db, search)
	if err != nil {
		return nil, err
	}

	var responses []dto.BarangResponse
	for _, barang := range barangList {
		responses = append(responses, *s.mapToResponse(&barang))
	}

	return &dto.BarangListResponse{
		Barang: responses,
		Total:  total,
	}, nil
}

func (s *barangServiceImpl) UpdateBarang(id int, req dto.BarangFormRequest) error {
	// Validasi dan konversi harga
	harga, err := strconv.ParseFloat(req.Harga, 64)
	if err != nil {
		return fmt.Errorf("format harga tidak valid")
	}

	if harga < 0 {
		return fmt.Errorf("harga tidak boleh negatif")
	}

	// Check if barang exists
	existingBarang, err := s.repo.GetByID(s.db, id)
	if err != nil {
		return err
	}

	var gambarPath = existingBarang.Gambar

	// Handle new image upload
	if req.Gambar != nil {
		newFilename, err := s.saveUploadedFile(req.Gambar)
		if err != nil {
			return fmt.Errorf("gagal upload gambar: %v", err)
		}

		// Delete old image if exists
		if existingBarang.Gambar != "" {
			oldPath := filepath.Join("uploads", "barang", existingBarang.Gambar)
			os.Remove(oldPath) // Ignore error if file doesn't exist
		}

		gambarPath = newFilename
	}

	barang := &model.Barang{
		ID:             id,
		NamaBarang:     req.NamaBarang,
		Deskripsi:      req.Deskripsi,
		Gambar:         gambarPath,
		Harga:          harga,
		LinkShopee:     req.LinkShopee,
		LinkTiktokshop: req.LinkTiktokshop,
	}

	return s.repo.Update(s.db, barang)
}

func (s *barangServiceImpl) DeleteBarang(id int) error {
	// Get barang to delete associated image
	barang, err := s.repo.GetByID(s.db, id)
	if err != nil {
		return err
	}

	// Delete from database first
	if err := s.repo.Delete(s.db, id); err != nil {
		return err
	}

	// Delete associated image file
	if barang.Gambar != "" {
		imagePath := filepath.Join("uploads", "barang", barang.Gambar)
		os.Remove(imagePath) // Ignore error if file doesn't exist
	}

	return nil
}

func (s *barangServiceImpl) saveUploadedFile(fileHeader *multipart.FileHeader) (string, error) {
	// Validate file type
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	isValidExt := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			isValidExt = true
			break
		}
	}

	if !isValidExt {
		return "", fmt.Errorf("tipe file tidak diizinkan. Gunakan: %s", strings.Join(allowedExtensions, ", "))
	}

	// Validate file size (max 5MB)
	if fileHeader.Size > 5*1024*1024 {
		return "", fmt.Errorf("ukuran file terlalu besar. Maksimal 5MB")
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, fileHeader.Filename)

	// Create directory if not exists
	uploadDir := filepath.Join("uploads", "barang")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("gagal membuat direktori upload: %v", err)
	}

	// Open source file
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file: %v", err)
	}
	defer src.Close()

	// Create destination file
	dstPath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("gagal membuat file tujuan: %v", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("gagal menyalin file: %v", err)
	}

	return filename, nil
}

func (s *barangServiceImpl) mapToResponse(barang *model.Barang) *dto.BarangResponse {
	response := &dto.BarangResponse{
		ID:             barang.ID,
		NamaBarang:     barang.NamaBarang,
		Deskripsi:      barang.Deskripsi,
		Gambar:         barang.Gambar,
		Harga:          barang.Harga,
		HargaFormatted: s.formatCurrency(barang.Harga),
		LinkShopee:     barang.LinkShopee,
		LinkTiktokshop: barang.LinkTiktokshop,
		CreatedAt:      barang.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      barang.UpdatedAt.Format(time.RFC3339),
	}

	// Generate full URL for image
	if barang.Gambar != "" {
		response.GambarURL = "/uploads/barang/" + barang.Gambar
	}

	return response
}

func (s *barangServiceImpl) formatCurrency(amount float64) string {
    // Konversi ke integer (bulatkan ke bawah)
    intAmount := int64(amount)
    
    // Konversi ke string
    str := fmt.Sprintf("%d", intAmount)

    var formatted []string
    for i, digit := range str {
        if i > 0 && (len(str)-i)%3 == 0 {
            formatted = append(formatted, ".")
        }
        formatted = append(formatted, string(digit))
    }
    
    return fmt.Sprintf("Rp. %s", strings.Join(formatted, ""))
}
