package models

import "time"

type VigenciaRango struct {
	FechaInicio time.Time `gorm:"column:fecha_inicio;not null" json:"fecha_inicio"`
	FechaFin    time.Time `gorm:"column:fecha_fin;not null" json:"fecha_fin"`
	IDPaquete   uint      `gorm:"column:id_paquete;not null" json:"id_paquete"`

	Paquete PaqueteTuristico `gorm:"foreignKey:IDPaquete;references:ID" json:"-"`
}

func (VigenciaRango) TableName() string {
	return "vigencias_rangos"
}
