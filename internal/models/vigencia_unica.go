package models

import "time"

type VigenciaUnica struct {
	Fecha     time.Time `gorm:"column:fecha;not null" json:"fecha"`
	IDPaquete uint      `gorm:"column:id_paquete;not null" json:"id_paquete"`

	Paquete PaqueteTuristico `gorm:"foreignKey:IDPaquete;references:ID" json:"-"`
}

func (VigenciaUnica) TableName() string {
	return "vigencias_unicas"
}
