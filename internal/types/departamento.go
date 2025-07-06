package types

import "time"

type DepartamentoTODO struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}
