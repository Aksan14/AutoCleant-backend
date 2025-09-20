package model

import "time"

type Barang struct {
	ID             int       `json:"id"`
	NamaBarang     string    `json:"nama_barang"`
	Deskripsi      string    `json:"deskripsi"`
	Gambar         string    `json:"gambar"`
	Harga          float64   `json:"harga"`
	LinkShopee     string    `json:"link_shopee"`
	LinkTiktokshop string    `json:"link_tiktokshop"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
