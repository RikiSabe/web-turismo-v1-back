package models

type EncargadoQR struct {
	ID   uint   `gorm:"column:id_foto;PrimaryKey;autoIncrement" json:"id"`
	Foto string `gorm:"column:foto" json:"foto"`

	IdUsuario uint `gorm:"column:id_usuario" json:"id_usuario"`

	Usuario Usuario `gorm:"foreignKey:IdUsuario;references:ID" json:"-"`
}

func (EncargadoQR) TableName() string {
	return "encargado_qr"
}
