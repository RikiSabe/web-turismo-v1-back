package types

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
	ID              uint      `json:"id"`
	Categoria       string    `json:"categoria"`
	Nombre          string    `json:"nombre"`
	Direccion       string    `json:"direccion"`
	Descripcion     string    `json:"descripcion"`
	HorarioApertura string    `json:"horario_apertura"`
	HorarioCierre   string    `json:"horario_cierre"`
	Precio          float64   `json:"precio"`
	Estado          bool      `json:"estado"`
	Encargado       Encargado `json:"encargado"`
	Ubicacion       Ubicacion `json:"ubicacion"`
	Fotos           any       `json:"fotos"`
}

type AtraccionTuristicaTODO struct {
	ID              uint      `json:"id"`
	Categoria       string    `json:"categoria"`
	Nombre          string    `json:"nombre"`
	Direccion       string    `json:"direccion"`
	Descripcion     string    `json:"descripcion"`
	HorarioApertura string    `json:"horario_apertura"`
	HorarioCierre   string    `json:"horario_cierre"`
	Precio          float64   `json:"precio"`
	Estado          bool      `json:"estado"`
	Encargado       Encargado `json:"encargado"`
	Ubicacion       Ubicacion `json:"ubicacion"`
}
