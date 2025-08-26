package models

import "time"

type Departamento struct {
	ID     uint   `gorm:"column:id_departamento;primaryKey;autoIncrement" json:"id"`
	Nombre string `gorm:"column:nombre;size:255;not null" json:"nombre"`
	Estado bool   `gorm:"column:estado;not null;default:true" json:"estado"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Departamento) TableName() string {
	return "departamentos"
}
