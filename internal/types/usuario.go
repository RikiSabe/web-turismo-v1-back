package types

type UsuarioTODO struct {
	ID        uint   `json:"id"`
	Rol       string `json:"rol"`
	Nombre    string `json:"nombre"`
	Apellido  string `json:"apellido"`
	Correo    string `json:"correo"`
	Telefono  string `json:"telefono"`
	Direccion string `json:"direccion"`
	Contra    string `json:"contra"`
	Estado    bool   `json:"estado"`
	Foto      string `json:"foto"`
}
