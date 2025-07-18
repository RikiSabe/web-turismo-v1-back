package types

import "time"

type UsuarioUnique struct {
	Rol             string `json:"rol"`
	Nombre          string `json:"nombre"`
	ApellidoPaterno string `json:"apellido_paterno"`
	ApellidoMaterno string `json:"apellido_materno"`
	FechaNacimiento string `json:"fecha_nacimiento"`
	CI              string `json:"ci"`
	Correo          string `json:"correo"`
	Telefono        string `json:"telefono"`
	IdUbicacion     uint   `json:"id_ubicacion"`
	Contra          string `json:"contra"`
	Estado          bool   `json:"estado"`
	Foto            string `json:"foto"`
}

type UsuarioResumen struct {
	Id              uint      `json:"id"`
	Rol             string    `json:"rol"`
	Nombre          string    `json:"nombre"`
	ApellidoPaterno string    `json:"apellido_paterno"`
	ApellidoMaterno string    `json:"apellido_materno"`
	FechaNacimiento time.Time `json:"fecha_nacimiento"`
	CI              string    `json:"ci"`
	Correo          string    `json:"correo"`
	Telefono        string    `json:"telefono"`
	IdUbicacion     string    `json:"id_ubicacion"`
	Estado          bool      `json:"estado"`
}

type UsuarioMenu struct {
	NombreCompleto string `json:"nombre_completo"`
	Foto           string `json:"foto"`
}
