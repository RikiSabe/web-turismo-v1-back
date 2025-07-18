package types

type AgenciaTODO struct {
	ID        uint   `json:"id"`
	Nombre    string `json:"nombre"`
	Direccion string `json:"direccion"`
	Telefono  string `json:"telefono"`
	Correo    string `json:"correo"`
	Estado    bool   `json:"estado"`
}
