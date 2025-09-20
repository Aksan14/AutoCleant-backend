package dto

type LoginRequest struct {
	NRA      string `json:"email"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	IdUser   string `json:"id_user" validate:"required"`
	NRA      string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ChangePasswordRequest struct {
	NRA              string `json:"email"`          
	OldPassword      string `json:"old_password"`
	NewPassword      string `json:"new_password"`
	ConfirmPassword  string `json:"confirm_password"`
}