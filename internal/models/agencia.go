package models

import (
	"time"
)

type Agencia struct {
	ID          uint   `gorm:"column:id_agencia;primaryKey;autoIncrement" json:"id"`
	Nombre      string `gorm:"column:nombre;size:255;not null" json:"nombre"`
	Descripcion string `gorm:"column:descripcion;size:255" json:"descripcion"`
	Direccion   string `gorm:"column:direccion;size:255;not null" json:"direccion"`
	Telefono    string `gorm:"column:telefono;size:50" json:"telefono"`
	Correo      string `gorm:"column:correo;size:255" json:"correo"`
	Estado      bool   `gorm:"column:estado;not null;default:true" json:"estado"`

	IdEncargado    uint `gorm:"column:id_encargado" json:"id_encargado"`
	IdDepartamento uint `gorm:"column:id_departamento" json:"id_departamento"`

	Usuario      Usuario      `gorm:"foreignKey:IdEncargado;references:ID" json:"-"`
	Departamento Departamento `gorm:"foreignKey:IdDepartamento;references:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Agencia) TableName() string {
	return "GestAgencias"
}
