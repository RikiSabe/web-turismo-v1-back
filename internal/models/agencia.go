package models

import (
	"time"
)

type Agencia struct {
	ID                uint   `gorm:"column:id_agencia;primaryKey;autoIncrement" json:"id"`
	Nombre            string `gorm:"size:255;not null" json:"nombre"`
	Direccion         string `gorm:"size:255;not null" json:"direccion"`
	Telefono          string `gorm:"size:50" json:"telefono"`
	CorreoElectronico string `gorm:"column:correo_electronico;size:255" json:"correo_electronico"`
	Estado            bool   `gorm:"not null;default:true" json:"estado"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Agencia) TableName() string {
	return "GestAgencias"
}
