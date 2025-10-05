package models

type AtraccionesCategorias struct {
	ID             uint `gorm:"primaryKey;autoIncrement" json:"id"`
	IDAtraccion    uint `gorm:"column:id_atraccion;not null" json:"id_atraccion"`
	IDSubCategoria uint `gorm:"column:id_subcategoria;not null" json:"id_subcategoria"`

	Atraccion    AtraccionTuristica `gorm:"foreignKey:IDAtraccion;references:ID" json:"-"`
	SubCategoria SubCategoria       `gorm:"foreignKey:IDSubCategoria;references:ID" json:"-"`
}

func TableName() string {
	return "atracciones_categorias"
}
