package types

type PaqueteTuristicoTODO struct {
	ID          uint    `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	Precio      float64 `json:"precio" format:"%.2f"`
	HoraInicial string  `json:"hora_inicial"`
	Estado      bool    `json:"estado"`
	Visible     bool    `json:"visible"`
}
