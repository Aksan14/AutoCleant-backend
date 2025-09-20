package dto

import "mime/multipart"

type BarangRequest struct {
	NamaBarang     string  `json:"nama_barang" validate:"required,min=1,max=255"`
	Deskripsi      string  `json:"deskripsi"`
	Harga          float64 `json:"harga" validate:"required,min=0"`
	LinkShopee     string  `json:"link_shopee" validate:"omitempty,url"`
	LinkTiktokshop string  `json:"link_tiktokshop" validate:"omitempty,url"`
}

type BarangFormRequest struct {
	NamaBarang     string                `form:"nama_barang" validate:"required,min=1,max=255"`
	Deskripsi      string                `form:"deskripsi"`
	Harga          string                `form:"harga" validate:"required"`
	LinkShopee     string                `form:"link_shopee" validate:"omitempty,url"`
	LinkTiktokshop string                `form:"link_tiktokshop" validate:"omitempty,url"`
	Gambar         *multipart.FileHeader `form:"gambar"`
}

type BarangResponse struct {
	ID             int     `json:"id"`
	NamaBarang     string  `json:"nama_barang"`
	Deskripsi      string  `json:"deskripsi"`
	Gambar         string  `json:"gambar"`
	GambarURL      string  `json:"gambar_url"`
	Harga          float64 `json:"harga"`
	HargaFormatted string  `json:"harga_formatted"`
	LinkShopee     string  `json:"link_shopee"`
	LinkTiktokshop string  `json:"link_tiktokshop"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type BarangListResponse struct {
	Barang []BarangResponse `json:"barang"`
	Total  int              `json:"total"`
}
