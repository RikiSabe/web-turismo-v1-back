package models

import (
	"time"
)

type Usuario struct {
	ID              uint      `gorm:"column:id_usuario;primaryKey;autoIncrement" json:"id"`
	Rol             string    `gorm:"column:rol;not null" json:"rol"`
	Nombre          string    `gorm:"column:nombre;size:255;not null" json:"nombre"`
	ApellidoPaterno string    `gorm:"column:apellido_paterno;size:255;not null" json:"apellido_paterno"`
	ApellidoMaterno string    `gorm:"column:apellido_materno;size:255" json:"apellido_materno"`
	FechaNacimiento time.Time `gorm:"column:fecha_nacimiento;not null" json:"fecha_nacimiento"`
	CI              string    `gorm:"column:ci;not null" json:"ci"`
	Correo          string    `gorm:"column:correo;size:255;not null" json:"correo"`
	Telefono        string    `gorm:"column:telefono;size:50" json:"telefono"`
	Contra          string    `gorm:"column:contra;size:255;not null" json:"contra"`
	Estado          bool      `gorm:"column:estado;not null;default:true" json:"estado"`
	Foto            string    `gorm:"column:foto;size:255" json:"foto"`
	IdUbicacion     uint      `gorm:"column:id_ubicacion" json:"id_ubicacion"`

	Provincia Provincia `gorm:"foreignKey:IdUbicacion;references:ID" json:"-"`
	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Usuario) TableName() string {
	return "GestUsuarios"
}
