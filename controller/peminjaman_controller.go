package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PeminjamanController interface {
	GetBarangTersedia(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	CreatePeminjaman(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	ReturnPeminjaman(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	ListPeminjaman(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	DeletePeminjaman(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}