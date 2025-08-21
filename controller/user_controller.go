package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type UserController interface {
	CreateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	LoginUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	FindByNRA(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	ChangePassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}