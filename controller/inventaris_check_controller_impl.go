package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"reset/dto"
	"reset/service"
)

type Controller interface {
	StartReport(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	AddCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdateCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	DeleteCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	FinalizeReport(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetReports(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetReportDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	ExportPDF(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type reportController struct {
	svc service.Service
}

func NewController(svc service.Service) Controller {
	return &reportController{svc: svc}
}

func (c *reportController) StartReport(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req dto.StartReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := c.svc.StartReport(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    resp,
		Message: "Report started successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *reportController) AddCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	reportID, _ := strconv.Atoi(idStr)

	var req dto.CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := c.svc.AddCheck(reportID, req)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
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
		Message: "Tambah Check Berhasil",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *reportController) UpdateCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, _ := strconv.Atoi(idStr)

	var req dto.CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.svc.UpdateCheck(id, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    nil,
		Message: "Check updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *reportController) DeleteCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, _ := strconv.Atoi(idStr)

	if err := c.svc.DeleteCheck(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    nil,
		Message: "Check deleted successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *reportController) FinalizeReport(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, _ := strconv.Atoi(idStr)

	if err := c.svc.FinalizeReport(id); err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    nil,
		Message: "Report finalized successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *reportController) GetReports(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	status := r.URL.Query().Get("status")

	reports, err := c.svc.GetReports(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    reports,
		Message: "Reports retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *reportController) GetReportDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, _ := strconv.Atoi(idStr)

	detail, err := c.svc.GetReportDetail(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    detail,
		Message: "Report detail retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *reportController) ExportPDF(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, _ := strconv.Atoi(idStr)

	data, err := c.svc.ExportPDF(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=report-%d.pdf", id))
	w.Write(data)
}