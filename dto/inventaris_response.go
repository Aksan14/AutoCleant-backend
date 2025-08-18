package dto

type ItemResponse struct {
    ID         int       `json:"id"`
    NamaBarang string    `json:"nama_barang"`
    Kategori   string    `json:"kategori"`
    Jumlah     int       `json:"jumlah"`
    Satuan     string    `json:"satuan"`
    Kondisi    string    `json:"kondisi"`
    // PhotoURL   string    `json:"photo_url"`
}
