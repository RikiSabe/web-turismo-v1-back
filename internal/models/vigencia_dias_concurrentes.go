package models

type VigenciaDiasConcurrentes struct {
	ID        uint `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDDia     uint `gorm:"column:id_dia;not null" json:"id_dia"`
	IDPaquete uint `gorm:"column:id_paquete;not null" json:"id_paquete"`

	Dia     Dia              `gorm:"foreignKey:IDDia;references:ID" json:"-"`
	Paquete PaqueteTuristico `gorm:"foreignKey:IDPaquete;references:ID" json:"-"`
}

func (VigenciaDiasConcurrentes) TableName() string {
	return "vigencias_dias_concurrentes"
}
