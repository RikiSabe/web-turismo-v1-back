package types

type Categoria struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Estado      bool   `json:"estado"`
}
