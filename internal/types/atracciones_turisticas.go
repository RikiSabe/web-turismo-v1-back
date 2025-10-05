package types

import (
	"gorm.io/datatypes"
)

type Foto struct {
	ID    uint   `json:"id"`
	Foto  string `json:"foto"`
	Orden uint   `json:"orden"`
}

type Encargado struct {
	ID             uint   `json:"id"`
	NombreCompleto string `json:"nombre_completo"`
}

type Ubicacion struct {
	Provincia struct {
		ID     uint   `json:"id"`
		Nombre string `json:"nombre"`
	} `json:"provincia"`
	Departamento struct {
		ID     uint   `json:"id"`
		Nombre string `json:"nombre"`
	} `json:"departamento"`
}

type AtraccionTuristicaUnique struct {
	ID              uint                          `json:"id"`
	Categoria       string                        `json:"categoria"`
	Nombre          string                        `json:"nombre"`
	Direccion       string                        `json:"direccion"`
	Descripcion     string                        `json:"descripcion"`
	HorarioApertura string                        `json:"horario_apertura"`
	HorarioCierre   string                        `json:"horario_cierre"`
	Precio          float64                       `json:"precio"`
	Estado          bool                          `json:"estado"`
	Encargado       datatypes.JSONType[Encargado] `json:"encargado"`
	Ubicacion       datatypes.JSONType[Ubicacion] `json:"ubicacion"`
	Fotos           datatypes.JSONType[[]Foto]    `json:"fotos"`
}

type AtraccionTuristicaTODO struct {
	ID              uint                          `json:"id"`
	Categoria       string                        `json:"categoria"`
	Nombre          string                        `json:"nombre"`
	Direccion       string                        `json:"direccion"`
	Descripcion     string                        `json:"descripcion"`
	HorarioApertura string                        `json:"horario_apertura"`
	HorarioCierre   string                        `json:"horario_cierre"`
	Precio          float64                       `json:"precio"`
	Estado          bool                          `json:"estado"`
	Encargado       datatypes.JSONType[Encargado] `json:"encargado"`
	Ubicacion       datatypes.JSONType[Ubicacion] `json:"ubicacion"`
}

type UsuarioEncargado struct {
	SRC            string `json:"src"`
	NombreCompleto string `json:"nombre_completo"`
	CI             string `json:"ci"`
	Correo         string `json:"correo"`
	Telefono       string `json:"telefono"`
}

type AtraccionesFotos struct {
	Nombre      string `json:"nombre"`
	Direccion   string `json:"direccion"`
	Descripcion string `json:"descripcion"`
	Foto        string `json:"foto"`
}
