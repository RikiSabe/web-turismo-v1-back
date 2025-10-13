package models

type Comprobante struct {
	ID        uint   `gorm:"column:id_comprobante;primaryKey;autoIncrement" json:"id"`
	IDReserva uint   `gorm:"column:id_reserva;not null" json:"id_reserva"`
	IDUsuario uint   `gorm:"column:id_usuario;not null" json:"id_usuario"`
	Foto      string `gorm:"column:foto" json:"foto"`

	Reserva Reservas `gorm:"foreignKey:IDReserva;references:ID" json:"-"`
	Usuario Usuario  `gorm:"foreignKey:IDUsuario;references:ID" json:"-"`
}

func (Comprobante) TableName() string {
	return "comprobantes"
}
