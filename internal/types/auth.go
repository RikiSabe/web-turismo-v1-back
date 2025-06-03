package types

type Usuario struct {
	Correo string `json:"correo"`
	Contra string `json:"contra"`
}

type RespuestaUsuario struct {
	ID  uint   `json:"id"`
	Rol string `json:"rol"`
}
