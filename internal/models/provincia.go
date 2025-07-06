package models

import "time"

type Provincia struct {
	ID     uint   `gorm:"column:id_provincia;primaryKey;autoIncrement" json:"id"`
	Nombre string `gorm:"size:255;not null" json:"nombre"`
	Estado bool   `gorm:"not null;default:true" json:"estado"`

	IDDepartamento uint `gorm:"column:id_departamento;not null" json:"id_departamento"`

	Departamento Departamento `gorm:"foreignKey:IDDepartamento;references:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Provincia) TableName() string {
	return "provincias"
}
