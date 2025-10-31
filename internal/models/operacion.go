package models

import "time"

type Operacion struct {
	ID               uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Origen           string    `gorm:"column:origen" json:"origen"`
	PersonasAdultas  uint      `gorm:"column:personas_adultas" json:"personas_adultas"`
	PersonasMenores  uint      `gorm:"column:personas_menores" json:"personas_menores"`
	PersonasNormales uint      `gorm:"column:personas_normales" json:"personas_normales"`
	Fecha            time.Time `gorm:"column:fecha" json:"fecha"`
	Precio           float64   `gorm:"column:precio" json:"precio"`
	RazonSocial      string    `gorm:"column:razon_social" json:"razon_social"`
	FormaPago        string    `gorm:"column:forma_pago" json:"forma_pago"`
	Identidad        string    `gorm:"column:identidad" json:"identidad"`

	IdPaquete   uint  `gorm:"column:id_paquete" json:"id_paquete"`
	IdEncargado *uint `gorm:"column:id_encargado" json:"id_encargado"`
	IdUbicacion uint  `gorm:"column:id_ubicacion" json:"id_ubicacion"`

	Usuario          Usuario          `gorm:"foreignKey:IdEncargado;references:ID" json:"-"`
	PaqueteTuristico PaqueteTuristico `gorm:"foreignKey:IdPaquete;references:ID" json:"-"`
	Ubicacion        Provincia        `gorm:"foreignKey:IdUbicacion;references:ID" json:"-"`
}

func (Operacion) TableName() string {
	return "operaciones"
}
