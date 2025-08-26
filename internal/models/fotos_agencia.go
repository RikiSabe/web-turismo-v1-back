package models

type FotosAgencia struct {
	ID        uint   `gorm:"column:id_foto;PrimaryKey;autoIncrement" json:"id"`
	IdAgencia uint   `gorm:"column:id_agencia" json:"id_agencia"`
	Foto      string `gorm:"column:foto;size:255" json:"foto"`
	Orden     uint   `gorm:"column:orden;orden" json:"orden"`

	Agencia Agencia `gorm:"foreignKey:IdAgencia;references:ID" json:"-"`
}

func (FotosAgencia) TableName() string {
	return "fotos_agencias"
}
