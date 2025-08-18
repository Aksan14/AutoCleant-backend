package dto

type PeminjamanResponse struct {
	ID             int     `json:"id"`
	BarangID       int     `json:"barang_id"`
	BarangNama     string  `json:"barang_nama"`
	NamaPeminjam   string  `json:"nama_peminjam"`
	TglPinjam      string  `json:"tgl_pinjam"`
	RencanaKembali string  `json:"rencana_kembali"`
	TglKembali     *string `json:"tgl_kembali,omitempty"`
	KondisiSetelah *string `json:"kondisi_setelah,omitempty"`
	Status         string  `json:"status"` // dipinjam | selesai
	Keterangan	 	string `json:"keterangan,omitempty"` // optional
}

type BarangSimpleResponse struct {
	ID         int    `json:"id"`
	NamaBarang string `json:"nama_barang"`
	Kategori   string `json:"kategori"`
	Satuan     string `json:"satuan"`
	Kondisi    string `json:"kondisi"`
	Foto       string `json:"foto"`
}