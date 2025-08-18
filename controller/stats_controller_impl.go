package controller

import (
	"net/http"
	"reset/dto"
	"reset/service"
	"reset/util"

	"github.com/julienschmidt/httprouter"
)

type reportControllerImpl struct {
	svc service.ReportService
}

// NewReportController membuat instance baru ReportController
func NewReportController(s service.ReportService) ReportController {
	return &reportControllerImpl{svc: s}
}

func (c *reportControllerImpl) CountAllBarang(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	count, err := c.svc.GetCountAllBarang(r.Context())
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
		Data:    count,
		Message: "Total barang berhasil diambil",
	}
	util.WriteJSON(w, http.StatusOK, response)
}

func (c *reportControllerImpl) CountBarangDipinjam(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	count, err := c.svc.GetCountBarangDipinjam(r.Context())
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
		Data:    count,
		Message: "Total barang dipinjam berhasil diambil",
	}
	util.WriteJSON(w, http.StatusOK, response)
}

func (c *reportControllerImpl) CountBarangRusakBerat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	count, err := c.svc.GetCountBarangRusakBerat(r.Context())
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
		Data:    count,
		Message: "Total barang rusak berhasil diambil",
	}
	util.WriteJSON(w, http.StatusOK, response)
}
