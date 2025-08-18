package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type InventarisController interface {
	CreateInventaris(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetByIDInventaris(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetAllInventaris(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdateInventaris(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	DeleteInventaris(w http.ResponseWriter, r *http.Request , ps httprouter.Params)
	SearchInventaris(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}
