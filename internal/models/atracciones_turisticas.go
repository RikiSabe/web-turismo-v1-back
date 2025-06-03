package models

import (
	"time"
)

type AtraccionTuristica struct {
	ID          uint    `gorm:"column:id_atracciones;primaryKey;autoIncrement" json:"id"`
	Tipo        string  `gorm:"size:255;not null" json:"tipo"`
	Nombre      string  `gorm:"size:255;not null" json:"nombre"`
	Ubicacion   string  `gorm:"size:255;not null" json:"ubicacion"`
	Descripcion string  `gorm:"type:text;not null" json:"descripcion"`
	Horarios    string  `gorm:"size:255" json:"horarios"`
	Precio      float64 `gorm:"type:numeric(10,2);not null" json:"precio"`
	Estado      bool    `gorm:"not null;default:true" json:"estado"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (AtraccionTuristica) TableName() string {
	return "GestAtraccionesTuristicas"
}
