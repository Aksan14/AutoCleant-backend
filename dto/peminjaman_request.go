package dto

type CreatePeminjamanRequest struct {
	BarangID       int    `json:"barang_id"`
	NamaPeminjam   string `json:"nama_peminjam"`
	TglPinjam      string `json:"tgl_pinjam"`       // "2006-01-02"
	RencanaKembali string `json:"rencana_kembali"`
	Keterangan     string `json:"keterangan,omitempty"` 
}

type ReturnPeminjamanRequest struct {
	TglKembali     string `json:"tgl_kembali"`      // "2006-01-02"
	KondisiSetelah string `json:"kondisi_setelah"`  // Baik/Rusak/dll
}