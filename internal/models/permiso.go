package models

import (
	"time"
)

type Permiso struct {
	ID     uint   `gorm:"column:id_permiso;primaryKey;autoIncrement" json:"id"`
	Nombre string `gorm:"size:255;not null" json:"nombre"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Permiso) TableName() string {
	return "GestPermisos"
}
