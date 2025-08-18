package repository

import (
	"context"
	"database/sql"
	"reset/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, tx *sql.Tx, user model.User) (model.User, error)
	FindById(ctx context.Context, tx *sql.Tx, idUser string) (model.User, error)
	FindByNRA(ctx context.Context, tx *sql.Tx, nra string) (model.User, error)
	UpdatePassword(ctx context.Context, tx *sql.Tx, user model.User) error
	CheckOldPassword(ctx context.Context, tx *sql.Tx, nra string) (string, error)
}