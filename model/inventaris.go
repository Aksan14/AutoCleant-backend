package model

import "time"

type Inventaris struct {
	ID         int       `json:"id"`
	NamaBarang string    `json:"Namabarang"`
	Kategori   string    `json:"Kategori"`
	Jumlah     int       `json:"Jumlah"`
	Satuan     string    `json:"Satuan"`
	Foto       string    `json:"Foto"`
	Kondisi    string    `json:"Kondisi"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Status     string		`json:"status"`
}