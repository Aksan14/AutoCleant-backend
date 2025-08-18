package util

import (
	"reset/dto"
	"reset/model"
)

func ConvertToResponseUsersDTO(user model.User) dto.UserResponse {
	return dto.UserResponse{
		IdUser: user.IdUser,
		NRA:    user.NRA,
	}
}