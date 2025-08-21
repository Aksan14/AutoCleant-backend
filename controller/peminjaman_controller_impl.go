package controller

import (
	"encoding/json"
	"net/http"
	"reset/dto"
	"reset/service"
	"reset/util"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type peminjamanControllerImpl struct {
	svc service.PeminjamanService
}

func (h *peminjamanControllerImpl) DeletePeminjaman(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "BAD_REQUEST",
			Message: "Invalid ID format",
		}
		util.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	err = h.svc.DeleteByID(r.Context(), id)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusInternalServerError,
			Status:  "ERROR",
			Message: err.Error(),
		}
		util.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "SUCCESS",
		Message: "Peminjaman berhasil dihapus",
	}
	util.WriteJSON(w, http.StatusOK, response)
}

func NewPeminjamanController(s service.PeminjamanService) PeminjamanController {
	return &peminjamanControllerImpl{svc: s}
}

func (h *peminjamanControllerImpl) GetBarangTersedia(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data, err := h.svc.ListBarangTersedia(r.Context())
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: err.Error(),
		}
		util.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	response := dto.ListResponseOK {
		Code: http.StatusOK,
		Status: "OK",
		Data: data,
		Message: "Barang tersedia berhasil diambil",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *peminjamanControllerImpl) CreatePeminjaman(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req dto.CreatePeminjamanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: err.Error(),
		}
		util.WriteJSON(w, http.StatusBadRequest, response)
		return
	}
	res, err := h.svc.CreatePeminjaman(r.Context(), req)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: err.Error(),
		}
		util.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	response := dto.ListResponseOK {
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    res,
		Message: "Peminjaman berhasil dibuat",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *peminjamanControllerImpl) ReturnPeminjaman(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/peminjaman/kembali/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: err.Error(),
		}
		util.WriteJSON(w, http.StatusBadRequest, response)
		return
	}
	var req dto.ReturnPeminjamanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: err.Error(),
		}
		util.WriteJSON(w, http.StatusBadRequest, response)
		return
	}
	if err := h.svc.ReturnPeminjaman(r.Context(), id, req); err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: err.Error(),
		}
		util.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    nil,
		Message: "Pengembalian berhasil",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *peminjamanControllerImpl) ListPeminjaman(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data, err := h.svc.ListPeminjaman(r.Context())
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: err.Error(),
		}
		util.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    data,
		Message: "Daftar peminjaman berhasil diambil",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}