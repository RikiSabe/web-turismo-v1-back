package models

import (
	"time"
)

type Usuario struct {
	ID        uint   `gorm:"column:id_usuario;primaryKey;autoIncrement" json:"id"`
	Rol       string `gorm:"not null" json:"rol"`
	Nombre    string `gorm:"size:255;not null" json:"nombre"`
	Apellido  string `gorm:"size:255;not null" json:"apellido"`
	Correo    string `gorm:"size:255;not null" json:"correo"`
	Telefono  string `gorm:"size:50" json:"telefono"`
	Direccion string `gorm:"size:255" json:"direccion"`
	Contra    string `gorm:"column:contra;size:255;not null" json:"contra"`
	Estado    bool   `gorm:"not null;default:true" json:"estado"`
	Foto      string `gorm:"size:255" json:"foto"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Usuario) TableName() string {
	return "GestUsuarios"
}
