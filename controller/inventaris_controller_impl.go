package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"reset/dto"
	"reset/model"
	"reset/service"

	"github.com/julienschmidt/httprouter"
)

type inventarisControllerImpl struct {
	inventarisService service.InventarisService
}

func NewInventarisController(service service.InventarisService) InventarisController {
	return &inventarisControllerImpl{inventarisService: service}
}

func saveUploadedFile(r *http.Request, formFileKey string) (string, error) {
	file, header, err := r.FormFile(formFileKey)
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", err
	}
	defer file.Close()

	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	filename := filepath.Join(uploadDir, header.Filename)
	out, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	cleanpath := filepath.Clean(filename)
	return cleanpath, nil
}

func (c *inventarisControllerImpl) CreateInventaris(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed parse multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	namaBarang := r.FormValue("Namabarang")
	kategori := r.FormValue("Kategori")
	jumlahStr := r.FormValue("Jumlah")
	satuan := r.FormValue("Satuan")
	kondisi := r.FormValue("Kondisi")

	jumlah, err := strconv.Atoi(jumlahStr)
	if err != nil {
		http.Error(w, "Invalid jumlah: "+err.Error(), http.StatusBadRequest)
		return
	}

	fotoPath, err := saveUploadedFile(r, "Foto")
	if err != nil {
		http.Error(w, "Failed save foto: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fotoPath = strings.ReplaceAll(fotoPath, "\\", "/")

	item := &model.Inventaris{
		NamaBarang: namaBarang,
		Kategori:   kategori,
		Jumlah:     jumlah,
		Satuan:     satuan,
		Kondisi:    kondisi,
		Foto:       fotoPath,
	}

	id, err := c.inventarisService.CreateInventaris(item)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Failed to create inventaris: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	item.ID = id
	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    item,
		Message: "Inventaris created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (c *inventarisControllerImpl) GetByIDInventaris(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	item, err := c.inventarisService.GetByIDInventaris(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if item == nil {
		http.NotFound(w, r)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    item,
		Message: "Inventaris retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *inventarisControllerImpl) GetAllInventaris(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	items, err := c.inventarisService.GetAllInventaris()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    items,
		Message: "Inventaris retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *inventarisControllerImpl) UpdateInventaris(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed parse multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	namaBarang := r.FormValue("Namabarang")
	kategori := r.FormValue("Kategori")
	jumlahStr := r.FormValue("Jumlah")
	satuan := r.FormValue("Satuan")
	kondisi := r.FormValue("Kondisi")

	jumlah, err := strconv.Atoi(jumlahStr)
	if err != nil {
		http.Error(w, "Invalid jumlah: "+err.Error(), http.StatusBadRequest)
		return
	}

	fotoPath, err := saveUploadedFile(r, "Foto")
	if err != nil {
		http.Error(w, "Failed save foto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fotoPath = strings.ReplaceAll(fotoPath, "\\", "/")

	item := &model.Inventaris{
		ID:         id,
		NamaBarang: namaBarang,
		Kategori:   kategori,
		Jumlah:     jumlah,
		Satuan:     satuan,
		Kondisi:    kondisi,
		Foto:       fotoPath,
	}

	if fotoPath == "" {
		existing, err := c.inventarisService.GetByIDInventaris(id)
		if err != nil {
			http.Error(w, "Data not found: "+err.Error(), http.StatusNotFound)
			return
		}
		if existing != nil {
			item.Foto = existing.Foto
		}
	}

	err = c.inventarisService.UpdateInventaris(item)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Failed to update inventaris: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    item,
		Message: "Inventaris updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *inventarisControllerImpl) DeleteInventaris(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid ID: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Ubah fungsi service.DeleteInventaris menjadi update kondisi "Dimusnahkan"
	err = c.inventarisService.MusnahkanInventaris(id)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Gagal memusnahkan inventaris: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    nil,
		Message: "Inventaris berhasil dimusnahkan",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *inventarisControllerImpl) SearchInventaris(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	keyword := r.URL.Query().Get("query")
	if keyword == "" {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Query parameter 'query' is required",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	items, err := c.inventarisService.SearchInventaris(keyword)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Error searching inventaris: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    items,
		Message: "Pencarian inventaris berhasil",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
