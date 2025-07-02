package types

import (
	"time"

	"gorm.io/datatypes"
)

type atraccionTuristicaEnlazada struct {
	ID          uint    `json:"id"`
	Tipo        string  `json:"tipo"`
	Nombre      string  `json:"nombre"`
	Ubicacion   string  `json:"ubicacion"`
	Descripcion string  `json:"descripcion"`
	Horarios    string  `json:"horarios"`
	Precio      float64 `json:"precio"`
	Estado      bool    `json:"estado"`
}

type PaqueteTuristicoTODO struct {
	ID          uint      `json:"id"`
	Categoria   string    `json:"categoria"`
	Nombre      string    `json:"nombre"`
	Fecha       time.Time `json:"fecha"`
	Descripcion string    `json:"descripcion"`
	Precio      float64   `json:"precio" format:"%.2f"`
	Duracion    string    `json:"duracion"`
	Salida      string    `json:"salida"`
	Estado      bool      `json:"estado"`

	IDAgencia uint `json:"id_agencia"`

	Atracciones datatypes.JSONType[[]atraccionTuristicaEnlazada] `json:"atracciones_turisticas"`
}
