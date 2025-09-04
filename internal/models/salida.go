package models

type Salida struct {
	ID   uint   `gorm:"column:id_salida;primaryKey;autoIncrement" json:"id"`
	Hora string `gorm:"column:hora;not null" json:"hora"`

	IDPaquete uint `gorm:"column:id_paquete;not null" json:"id_paquete"`

	Paquete PaqueteTuristico `gorm:"foreignKey:IDPaquete;references:ID" json:"-"`
}

func (Salida) TableName() string {
	return "salidas"
}
