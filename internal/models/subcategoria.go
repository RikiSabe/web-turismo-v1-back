package models

import "time"

type SubCategoria struct {
	ID          uint   `gorm:"column:id_subcategoria;primaryKey;autoIncrement" json:"id"`
	Nombre      string `gorm:"column:nombre;size:255;not null" json:"nombre"`
	Descripcion string `gorm:"column:descripcion;size:255;not null" json:"descripcion"`
	Estado      bool   `gorm:"column:estado;not null;default:true" json:"estado"`

	IDCategoria uint `gorm:"column:id_categoria;not null" json:"id_categoria"`

	Categoria Categoria `gorm:"foreignKey:IDCategoria;references:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (SubCategoria) TableName() string {
	return "subcategorias"
}
