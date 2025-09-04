package models

import (
	"time"
)

type PaqueteTuristico struct {
	ID          uint    `gorm:"column:id_paquete_turistico;primaryKey;autoIncrement" json:"id"`
	Nombre      string  `gorm:"column:nombre;size:255;not null" json:"nombre"`
	Descripcion string  `gorm:"column:descripcion;size:500;not null" json:"descripcion"`
	Actividades string  `gorm:"column:actividades;size:500;not null" json:"actividades"`
	Tipo        string  `gorm:"column:tipo;size:100;not null" json:"tipo"`
	Precio      float64 `gorm:"column:precio;not null" json:"precio" format:"%.2f"`
	Duracion    string  `gorm:"column:duracion;not null" json:"duracion"`
	HoraInicial string  `gorm:"column:hora_inicial;not null" json:"hora_inicial"`
	HoraFinal   string  `gorm:"column:hora_final;not null" json:"hora_final"`
	CupoTotal   uint    `gorm:"column:cupo_total;not null" json:"cupo_total"`
	CupoMinimo  uint    `gorm:"column:cupo_minimo;not null" json:"cupo_minimo"`

	Estado      bool `gorm:"column:estado;default:true" json:"estado"`
	Visible     bool `gorm:"column:visible;default:true" json:"visible"`
	Promocional bool `gorm:"column:promocional;default:false" json:"promocional"`

	IDAgencia uint `gorm:"column:id_agencia;not null" json:"id_agencia"`

	Agencia Agencia `gorm:"foreignKey:IDAgencia;references:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (PaqueteTuristico) TableName() string {
	return "GestPaquetesTuristicos"
}
