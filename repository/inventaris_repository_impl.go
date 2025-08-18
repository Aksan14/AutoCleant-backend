package repository

import "reset/model"

type ItemRepository interface {
    CreateInventaris(item *model.Inventaris) (int, error)
    GetByIDInventaris(id int) (*model.Inventaris, error)
    GetAllInventaris() ([]*model.Inventaris, error)
    UpdateInventaris(item *model.Inventaris) error
    DeleteInventaris(id int, kondisi string) error
    FindByNamabarang(namaBarang string) (*model.Inventaris, error)
    IsBarangSedangDipinjam(barangID int) (bool, error)
    SearchInventaris(keyword string) ([]model.Inventaris, error)
}