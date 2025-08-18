package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ReportController interface {
	CountAllBarang(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	CountBarangDipinjam(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	CountBarangRusakBerat(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
}
