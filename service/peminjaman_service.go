package service

import (
	"context"
	"reset/dto"
)

type PeminjamanService interface {
	ListBarangTersedia(ctx context.Context) ([]dto.BarangSimpleResponse, error)
	CreatePeminjaman(ctx context.Context, req dto.CreatePeminjamanRequest) (dto.PeminjamanResponse, error)
	ReturnPeminjaman(ctx context.Context, id int, req dto.ReturnPeminjamanRequest) error
	ListPeminjaman(ctx context.Context) ([]dto.PeminjamanResponse, error)
	DeleteByID(ctx context.Context, id int) error
}
