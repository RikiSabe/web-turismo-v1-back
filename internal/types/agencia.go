package types

type AgenciaTODO struct {
	ID                uint   `json:"id"`
	Nombre            string `json:"nombre"`
	Direccion         string `json:"direccion"`
	Telefono          string `json:"telefono"`
	CorreoElectronico string `json:"correo_electronico"`
	Estado            bool   `json:"estado"`
}
