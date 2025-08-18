package controller

import (
	"encoding/json"
	"net/http"
	"reset/dto"
	"reset/service"
	"reset/util"

	"github.com/julienschmidt/httprouter"
)

type userControllerImpl struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userControllerImpl{
		UserService: userService,
	}
}

func (controller *userControllerImpl) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var requestCreate dto.CreateUserRequest
	util.ReadFromRequestBody(r, &requestCreate)

	// Cek user sudah ada
	existingUser, _ := controller.UserService.FindByNRA(r.Context(), requestCreate.NRA)
	if existingUser != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "NRA sudah terdaftar",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		util.WriteToResponseBody(w, response)
		return
	}
	
	responseDTO := controller.UserService.CreateUser(r.Context(), requestCreate)
	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    responseDTO,
		Message: "User berhasil dibuat",
	}

	w.Header().Set("Content-Type", "application/json")
	util.WriteToResponseBody(w, response)
}

func (controller *userControllerImpl) LoginUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "FAILED",
			Message: "Invalid input",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		util.WriteToResponseBody(w, response)
		return
	}

	token, err := controller.UserService.LoginUser(r.Context(), loginRequest)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusUnauthorized,
			Status:  "FAILED",
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		util.WriteToResponseBody(w, response)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    token,
		Message: "Token generated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	util.WriteToResponseBody(w, response)
}

func (controller *userControllerImpl) ChangePassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "FAILED",
			Message: "Invalid input",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		util.WriteToResponseBody(w, response)
		return
	}

	err := controller.UserService.ChangePassword(r.Context(), req)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "FAILED",
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		util.WriteToResponseBody(w, response)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Password berhasil diubah",
	}

	w.Header().Set("Content-Type", "application/json")
	util.WriteToResponseBody(w, response)
}