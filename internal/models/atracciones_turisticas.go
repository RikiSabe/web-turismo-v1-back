package models

import (
	"time"
)

type AtraccionTuristica struct {
	ID              uint    `gorm:"column:id_atraccion;primaryKey;autoIncrement" json:"id"`
	IdEncargado     uint    `gorm:"column:id_encargado" json:"id_encargado"`
	IdUbicacion     uint    `gorm:"column:id_ubicacion" json:"id_ubicacion"`
	Categoria       string  `gorm:"column:categoria;size:255;not null" json:"categoria"` // eliminar
	Nombre          string  `gorm:"column:nombre;size:255;not null" json:"nombre"`
	Direccion       string  `gorm:"column:direccion;size:255;not null" json:"direccion"`
	Descripcion     string  `gorm:"column:descripcion;not null" json:"descripcion"`
	HorarioApertura string  `gorm:"column:horario_apertura;size:255" json:"horario_apertura"`
	HorarioCierre   string  `gorm:"column:horario_cierre;size:255" json:"horario_cierre"`
	Precio          float64 `gorm:"column:precio;type:numeric(10,2);not null" json:"precio"`
	Estado          bool    `gorm:"column:estado;not null;default:true" json:"estado"`

	Usuario   Usuario   `gorm:"foreignKey:IdEncargado;references:ID" json:"-"`
	Provincia Provincia `gorm:"foreignKey:IdUbicacion;references:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (AtraccionTuristica) TableName() string {
	return "GestAtraccionesTuristicas"
}
