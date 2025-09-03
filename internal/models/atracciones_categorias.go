package models

type AtraccionesCategorias struct {
	ID          uint `gorm:"primaryKey;autoIncrement" json:"id"`
	IDAtraccion uint `gorm:"column:id_atraccion;not null" json:"id_atraccion"`
	IDCategoria uint `gorm:"column:id_categoria;not null" json:"id_categoria"`

	Atraccion AtraccionTuristica `gorm:"foreignKey:IDAtraccion;references:ID" json:"-"`
	Categoria Categoria          `gorm:"foreignKey:IDCategoria;references:ID" json:"-"`
}

func TableName() string {
	return "atracciones_categorias"
}
