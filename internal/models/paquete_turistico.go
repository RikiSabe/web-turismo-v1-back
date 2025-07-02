package models

import (
	"time"
)

type PaqueteTuristico struct {
	ID          uint      `gorm:"column:id_paquete_turistico;primaryKey;autoIncrement" json:"id"`
	Categoria   string    `gorm:"size:255;not null" json:"categoria"`
	Nombre      string    `gorm:"size:255;not null" json:"nombre"`
	Fecha       time.Time `gorm:"not null" json:"fecha"`
	Descripcion string    `gorm:"size:500;not null" json:"descripcion"`
	Precio      float64   `gorm:"not null" json:"precio" format:"%.2f"`
	Duracion    string    `gorm:"not null" json:"duracion"`
	Salida      string    `gorm:"size:255;not null" json:"salida"`
	Estado      bool      `gorm:"size:50;not null" json:"estado"`

	IDAgencia uint `gorm:"column:id_agencia;not null" json:"id_agencia"`

	Agencia Agencia `gorm:"foreignKey:IDAgencia;references:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (PaqueteTuristico) TableName() string {
	return "GestPaquetesTuristicos"
}
