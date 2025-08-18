package repository

import (
	"context"
	"database/sql"
	"errors"
	"reset/model"
)

type userRepositoryImpl struct {
	DB *sql.DB
}

func NewUserRepositoryImpl(db *sql.DB) UserRepository {
	return &userRepositoryImpl{DB: db}
}

func (r *userRepositoryImpl) CreateUser(ctx context.Context, tx *sql.Tx, user model.User) (model.User, error) {
	query := `INSERT INTO users(id, nra, password) VALUES(?, ?, ?)`

	_, err := tx.ExecContext(ctx, query, user.IdUser, user.NRA, user.Password)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *userRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, idUser string) (model.User, error) {
	query := `SELECT id, nra, password FROM users WHERE id = ?`

	row := tx.QueryRowContext(ctx, query, idUser)
	user := model.User{}
	err := row.Scan(&user.IdUser, &user.NRA, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, err
	}

	return user, nil
}

func (r *userRepositoryImpl) FindByNRA(ctx context.Context, tx *sql.Tx, nra string) (model.User, error) {
	query := `SELECT id, nra, password FROM users WHERE nra = ?`

	row := tx.QueryRowContext(ctx, query, nra)
	user := model.User{}
	err := row.Scan(&user.IdUser, &user.NRA, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, err
	}

	return user, nil
}

func (r *userRepositoryImpl) UpdatePassword(ctx context.Context, tx *sql.Tx, user model.User) error {
	_, err := tx.ExecContext(ctx,
		"UPDATE users SET password=? WHERE nra=?",
		user.Password, user.NRA,
	)
	return err
}

func (r *userRepositoryImpl) CheckOldPassword(ctx context.Context, tx *sql.Tx, nra string) (string, error) {
	var hashedPassword string

	err := tx.QueryRowContext(ctx,
		"SELECT password FROM users WHERE nra = ?",
		nra,
	).Scan(&hashedPassword)

	if err != nil {
		return "", err
	}

	return hashedPassword, nil
}

