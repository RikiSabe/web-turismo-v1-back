package models

import "time"

type Reservas struct {
	ID             uint      `gorm:"column:id_reserva;primaryKey;autoIncrement" json:"id"`
	Fecha          time.Time `gorm:"not null" json:"fecha"`
	Descripcion    string    `gorm:"size:500;not null" json:"descripcion"`
	NumeroPersonas int       `gorm:"not null" json:"numero_personas"`
	Estado         bool      `gorm:"not null" json:"estado"`

	IDUsuario uint `gorm:"column:id_usuario;not null" json:"id_usuario"`
	IDPaquete uint `gorm:"column:id_paquete;not null" json:"id_paquete"`
	// Relaciones
	Paquete PaqueteTuristico `gorm:"foreignKey:IDPaquete;references:ID" json:"-"`
	Usuario Usuario          `gorm:"foreignKey:IDUsuario;references:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Reservas) TableName() string {
	return "GestReservas"
}
