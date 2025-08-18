package model

import "time"

type Peminjaman struct {
	ID              int        `json:"id"`
	BarangID        int        `json:"barang_id"`
	NamaPeminjam    string     `json:"nama_peminjam"`
	TglPinjam       time.Time  `json:"tgl_pinjam"`
	RencanaKembali  time.Time  `json:"rencana_kembali"`
	TglKembali      *time.Time `json:"tgl_kembali,omitempty"`
	KondisiSetelah  *string    `json:"kondisi_setelah,omitempty"`
	Status          string     `json:"status"` // dipinjam | selesai
	Keterangan	  	string    `json:"keterangan,omitempty"` // optional
	// optional join
	BarangNama string `json:"barang_nama,omitempty"`
}
