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

	existingUser, _ := controller.UserService.FindByNRA(r.Context(), requestCreate.NRA)
	if existingUser != nil {
		response := dto.ListResponseError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Email sudah terdaftar",
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

func (controller *userControllerImpl) FindByNRA(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	nra := ps.ByName("email")
	user, err := controller.UserService.FindByNRA(r.Context(), nra)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusNotFound,
			Status:  "FAILED",
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		util.WriteToResponseBody(w, response)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Data:    user,
		Message: "User found",
	}

	w.Header().Set("Content-Type", "application/json")
	util.WriteToResponseBody(w, response)
}

func (controller *userControllerImpl) ChangePassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var changePasswordRequest dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&changePasswordRequest); err != nil {
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

	err := controller.UserService.ChangePassword(r.Context(), changePasswordRequest)
	if err != nil {
		response := dto.ListResponseError{
			Code:    http.StatusInternalServerError,
			Status:  "FAILED",
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		util.WriteToResponseBody(w, response)
		return
	}

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Password changed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	util.WriteToResponseBody(w, response)
}
