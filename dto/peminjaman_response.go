package dto

type PeminjamanResponse struct {
	ID                int     `json:"id"`
	BarangID          int     `json:"barang_id"`
	BarangNama        string  `json:"barang_nama"`
	NamaPeminjam      string  `json:"nama_peminjam"`
	TglPinjam         string  `json:"tgl_pinjam"`
	RencanaKembali    string  `json:"rencana_kembali"`
	Jumlah            int     `json:"jumlah"`
	TglKembali        *string `json:"tgl_kembali,omitempty"`
	KondisiSetelah    *string `json:"kondisi_setelah,omitempty"`
	Status            string  `json:"status"`               // dipinjam | selesai
	Keterangan        string  `json:"keterangan,omitempty"` // optional
	FotoBuktiKembali  *string `json:"foto_bukti_kembali,omitempty"`
	KeteranganKembali *string `json:"keterangan_kembali,omitempty"`
}

type BarangSimpleResponse struct {
	ID         int    `json:"id"`
	NamaBarang string `json:"nama_barang"`
	Kategori   string `json:"kategori"`
	Satuan     string `json:"satuan"`
	Kondisi    string `json:"kondisi"`
	Foto       string `json:"foto"`
	Jumlah     int    `json:"jumlah"`
}
