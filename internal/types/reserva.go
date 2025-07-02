package types

import "time"

type ReservaTODO struct {
	ID             uint      `json:"id"`
	Fecha          time.Time `json:"fecha"`
	Descripcion    string    `json:"descripcion"`
	NumeroPersonas int       `json:"numero_personas"`
	Estado         bool      `json:"estado"`

	IDUsuario uint `json:"id_usuario"`
	IDPaquete uint `json:"id_paquete"`
}
