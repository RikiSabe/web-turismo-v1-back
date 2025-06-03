package models

import (
	"time"
)

type Rol struct {
	ID     uint   `gorm:"column:id_rol;primaryKey;autoIncrement" json:"id"`
	Nombre string `gorm:"size:255;not null" json:"nombre"`
	Estado bool   `gorm:"not null;default:true" json:"estado"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Rol) TableName() string {
	return "GestRoles"
}
