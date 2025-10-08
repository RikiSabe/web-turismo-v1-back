package models

type FotosAtracciones struct {
	ID          uint   `gorm:"column:id_foto;PrimaryKey;autoIncrement" json:"id"`
	IdAtraccion uint   `gorm:"column:id_atraccion" json:"id_atraccion"`
	Foto        string `gorm:"column:foto" json:"foto"`
	Orden       uint   `gorm:"column:orden;orden" json:"orden"`

	AtraccionTuristica AtraccionTuristica `gorm:"foreignKey:IdAtraccion;references:ID" json:"-"`
}

func (FotosAtracciones) TableName() string {
	return "fotos_atracciones"
}
