package dto

// For create/update we expect multipart/form-data for photo.
// Other fields can be read from form values.
type ItemRequest struct {
    NamaBarang string
    Kategori   string
    Jumlah     int
    Satuan     string
    Kondisi    string
    // Foto      string
    // Photo will be handled from multipart file; not in struct here.
}
