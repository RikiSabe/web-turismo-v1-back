package models

type PaqueteAtraccion struct {
	ID          uint `gorm:"primaryKey;autoIncrement" json:"id"`
	IDPaquete   uint `gorm:"column:id_paquete;not null" json:"id_paquete"`
	IDAtraccion uint `gorm:"column:id_atraccion;not null" json:"id_atraccion"`

	Paquete   PaqueteTuristico   `gorm:"foreignKey:IDPaquete;references:ID" json:"-"`
	Atraccion AtraccionTuristica `gorm:"foreignKey:IDAtraccion;references:ID" json:"-"`
}

func (PaqueteAtraccion) TableName() string {
	return "PaquetesAtracciones"
}
