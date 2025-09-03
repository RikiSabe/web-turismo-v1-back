package models

import "time"

type Categoria struct {
	ID          uint   `gorm:"column:id_categoria;primaryKey;autoIncrement" json:"id"`
	Nombre      string `gorm:"column:nombre;size:255;not null" json:"nombre"`
	Descripcion string `gorm:"column:descripcion;size:255;not null" json:"descripcion"`
	Estado      bool   `gorm:"column:estado;not null;default:true" json:"estado"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Categoria) TableName() string {
	return "categorias"
}
