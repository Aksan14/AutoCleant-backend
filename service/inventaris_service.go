package service

import "reset/model"

type InventarisService interface {
	CreateInventaris(item *model.Inventaris) (int, error)
	GetByIDInventaris(id int) (*model.Inventaris, error)
	GetAllInventaris() ([]*model.Inventaris, error)
	UpdateInventaris(item *model.Inventaris) error
	DeleteInventaris(item *model.Inventaris) error
	MusnahkanInventaris(id int) error
	SearchInventaris(keyword string) ([]model.Inventaris, error)
}