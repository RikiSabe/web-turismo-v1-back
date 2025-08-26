package types

import "gorm.io/datatypes"

type AgenciaTODO struct {
	ID             uint   `json:"id"`
	Nombre         string `json:"nombre"`
	Direccion      string `json:"direccion"`
	Telefono       string `json:"telefono"`
	Correo         string `json:"correo"`
	Estado         bool   `json:"estado"`
	IdEncargado    uint   `json:"id_encargado"`
	IdDepartamento uint   `json:"id_departamento"`
}

type AgenciaUnique struct {
	ID             uint                       `json:"id"`
	Nombre         string                     `json:"nombre"`
	Direccion      string                     `json:"direccion"`
	Telefono       string                     `json:"telefono"`
	Correo         string                     `json:"correo"`
	Descripcion    string                     `json:"descripcion"`
	Estado         bool                       `json:"estado"`
	IdEncargado    uint                       `json:"id_encargado"`
	IdDepartamento uint                       `json:"id_departamento"`
	Fotos          datatypes.JSONType[[]Foto] `json:"fotos"`
}
