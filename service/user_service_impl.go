package service

import (
	"context"
	"database/sql"
	"errors"
	"reset/dto"
	"reset/model"
	"reset/repository"
	"reset/util"

	"github.com/google/uuid"
)

type userServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
}

func NewUserServiceImpl(userRepository repository.UserRepository, db *sql.DB) UserService {
	return &userServiceImpl{
		UserRepository: userRepository,
		DB:             db,
	}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, req dto.CreateUserRequest) dto.UserResponse {
	tx, err := s.DB.BeginTx(ctx, nil)
	util.SentPanicIfError(err)
	defer util.CommitOrRollBack(tx)

	hashedPass, err := util.HashPassword(req.Password)
	util.SentPanicIfError(err)

	user := model.User{
		IdUser:   uuid.New().String(),
		NRA:      req.NRA,
		Password: hashedPass,
	}

	newUser, err := s.UserRepository.CreateUser(ctx, tx, user)
	util.SentPanicIfError(err)

	return util.ConvertToResponseUsersDTO(newUser)
}

func (s *userServiceImpl) LoginUser(ctx context.Context, req dto.LoginRequest) (string, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer util.CommitOrRollBack(tx)

	user, err := s.UserRepository.FindByNRA(ctx, tx, req.NRA)
	if err != nil || !util.VerifyPassword(user.Password, req.Password) {
		return "", errors.New("")
	}

	token, err := util.GenerateJWT(user.NRA)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userServiceImpl) FindByNRA(ctx context.Context, nra string) (*dto.UserResponse, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer util.CommitOrRollBack(tx)

	user, err := s.UserRepository.FindByNRA(ctx, tx, nra)
	if err != nil {
		return nil, err
	}

	response := util.ConvertToResponseUsersDTO(user)
	return &response, nil
}

func (s *userServiceImpl) ChangePassword(ctx context.Context, request dto.ChangePasswordRequest) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer util.CommitOrRollBack(tx)

	// Ambil password lama dari DB
	oldHashedPass, err := s.UserRepository.CheckOldPassword(ctx, tx, request.NRA)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	// Cek password lama
	if !util.VerifyPassword(oldHashedPass, request.OldPassword) {
		return errors.New("password lama salah")
	}

	// Cek konfirmasi password
	if request.NewPassword != request.ConfirmPassword {
		return errors.New("konfirmasi password tidak cocok")
	}

	// Hash password baru
	newHashedPass, err := util.HashPassword(request.NewPassword)
	if err != nil {
		return errors.New("gagal meng-hash password baru")
	}

	// Update password user
	user := model.User{
		NRA:      request.NRA,
		Password: newHashedPass,
	}

	if err := s.UserRepository.UpdatePassword(ctx, tx, user); err != nil {
		return errors.New("gagal update password")
	}

	return nil
}
