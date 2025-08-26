package models

import (
	"time"
)

type PaqueteTuristico struct {
	ID          uint      `gorm:"column:id_paquete_turistico;primaryKey;autoIncrement" json:"id"`
	Categoria   string    `gorm:"column:categoria;size:255;not null" json:"categoria"`
	Nombre      string    `gorm:"column:nombre;size:255;not null" json:"nombre"`
	Fecha       time.Time `gorm:"column:fecha;not null" json:"fecha"`
	Descripcion string    `gorm:"column:descripcion;size:500;not null" json:"descripcion"`
	Precio      float64   `gorm:"column:precio;not null" json:"precio" format:"%.2f"`
	Duracion    string    `gorm:"column:duracion;not null" json:"duracion"`
	Salida      string    `gorm:"column:salida;size:255;not null" json:"salida"`
	Estado      bool      `gorm:"column:estado;not null" json:"estado"`

	IDAgencia uint `gorm:"column:id_agencia;not null" json:"id_agencia"`

	Agencia Agencia `gorm:"foreignKey:IDAgencia;references:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (PaqueteTuristico) TableName() string {
	return "GestPaquetesTuristicos"
}
