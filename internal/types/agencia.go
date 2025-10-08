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

type AgenciaDatosGenerales struct {
	Nombre         string `json:"nombre"`
	Telefono       string `json:"telefono"`
	Correo         string `json:"correo"`
	IdEncargado    uint   `json:"id_encargado"`
	Descripcion    string `json:"descripcion"`
	IdDepartamento uint   `json:"id_departamento"`
	Direccion      string `json:"direccion"`
	Estado         bool   `json:"estado"`
}

type AgenciaFotos struct {
	Fotos datatypes.JSONType[[]Foto] `json:"fotos"`
}
