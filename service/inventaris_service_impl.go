package service

import (
	"errors"
	"reset/model"
	"reset/repository"
)

type inventarisServiceImpl struct {
	repo repository.ItemRepository
}

func NewInventarisService(repo repository.ItemRepository) InventarisService {
	return &inventarisServiceImpl{repo: repo}
}

// Create
func (s *inventarisServiceImpl) CreateInventaris(item *model.Inventaris) (int, error) {
	existing, err := s.repo.FindByNamabarang(item.NamaBarang)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, errors.New("nama barang sudah terdaftar")
	}

	return s.repo.CreateInventaris(item)
}

func (s *inventarisServiceImpl) UpdateInventaris(item *model.Inventaris) error {
	// Cek apakah ID ada
	oldData, err := s.repo.GetByIDInventaris(item.ID)
	if err != nil {
		return err
	}
	if oldData == nil {
		return errors.New("data tidak ditemukan")
	}

	// Cek duplikasi nama barang
	existing, err := s.repo.FindByNamabarang(item.NamaBarang)
	if err != nil {
		return err
	}
	if existing != nil && existing.ID != item.ID {
		return errors.New("nama barang sudah terdaftar")
	}

	return s.repo.UpdateInventaris(item)
}
// Delete

func (s *inventarisServiceImpl) GetByIDInventaris(id int) (*model.Inventaris, error) {
	return s.repo.GetByIDInventaris(id)
}

func (s *inventarisServiceImpl) GetAllInventaris() ([]*model.Inventaris, error) {
	return s.repo.GetAllInventaris()
}

func (s *inventarisServiceImpl) DeleteInventaris(item *model.Inventaris) error {
    existing, err := s.repo.GetByIDInventaris(item.ID)
    if err != nil {
        return err
    }

    if existing.Kondisi == "Dimusnahkan" {
        return errors.New("barang sudah dimusnahkan dan tidak dapat diupdate")
    }

    duplicate, err := s.repo.FindByNamabarang(item.NamaBarang)
    if err != nil {
        return err
    }
    if duplicate != nil && duplicate.ID != item.ID {
        return errors.New("nama barang sudah terdaftar")
    }

    return s.repo.UpdateInventaris(item)
}

func (s *inventarisServiceImpl) MusnahkanInventaris(id int) error {
    return s.repo.DeleteInventaris(id, "Dimusnahkan")
}

func (s *inventarisServiceImpl) SearchInventaris(keyword string) ([]model.Inventaris, error) {
    return s.repo.SearchInventaris(keyword)
}