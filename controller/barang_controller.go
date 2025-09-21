package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"reset/dto"
	"reset/service"

	"github.com/julienschmidt/httprouter"
)

type BarangController interface {
	CreateBarang(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	GetBarang(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetAllBarang(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	UpdateBarang(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	DeleteBarang(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type barangController struct {
	svc service.BarangService
}

func NewBarangController(svc service.BarangService) BarangController {
	return &barangController{svc: svc}
}

func (c *barangController) CreateBarang(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Failed to parse multipart form",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	req := dto.BarangFormRequest{
		NamaBarang:     r.FormValue("nama_barang"),
		Deskripsi:      r.FormValue("deskripsi"),
		Harga:          r.FormValue("harga"),
		LinkShopee:     r.FormValue("link_shopee"),
		LinkTiktokshop: r.FormValue("link_tiktokshop"),
	}

	// Validasi required fields
	if req.NamaBarang == "" {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Nama barang harus diisi",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.Harga == "" {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Harga harus diisi",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Handle file upload
	file, fileHeader, err := r.FormFile("gambar")
	if err == nil {
		defer file.Close()
		req.Gambar = fileHeader
	}

	resp, err := c.svc.CreateBarang(req)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    resp,
		Message: "Barang created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (c *barangController) GetBarang(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid ID format",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	resp, err := c.svc.GetBarangByID(id)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    resp,
		Message: "Barang retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *barangController) GetAllBarang(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	search := r.URL.Query().Get("search")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	resp, err := c.svc.GetAllBarang(search, page, limit)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    resp,
		Message: "Barang list retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *barangController) UpdateBarang(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid ID format",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Failed to parse multipart form",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	req := dto.BarangFormRequest{
		NamaBarang:     r.FormValue("nama_barang"),
		Deskripsi:      r.FormValue("deskripsi"),
		Harga:          r.FormValue("harga"),
		LinkShopee:     r.FormValue("link_shopee"),
		LinkTiktokshop: r.FormValue("link_tiktokshop"),
	}

	// Validasi required fields
	if req.NamaBarang == "" {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Nama barang harus diisi",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.Harga == "" {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Harga harus diisi",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Handle file upload
	file, fileHeader, err := r.FormFile("gambar")
	if err == nil {
		defer file.Close()
		req.Gambar = fileHeader
	}

	if err := c.svc.UpdateBarang(id, req); err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    nil,
		Message: "Barang updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *barangController) DeleteBarang(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid ID format",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := c.svc.DeleteBarang(id); err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    nil,
		Message: "Barang deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
