package models

type Dia struct {
	ID     uint   `gorm:"column:id_dia;primaryKey;autoIncrement" json:"id_dia"`
	Nombre string `gorm:"column:nombre;not null" json:"nombre"`
}

func (Dia) TableName() string {
	return "dias"
}
