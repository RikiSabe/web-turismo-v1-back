package types

type AtraccionTuristicaTODO struct {
	ID          uint    `json:"id"`
	Tipo        string  `json:"tipo"`
	Nombre      string  `json:"nombre"`
	Ubicacion   string  `json:"ubicacion"`
	Descripcion string  `json:"descripcion"`
	Horarios    string  `json:"horarios"`
	Precio      float64 `json:"precio"`
	Estado      bool    `json:"estado"`
}
