package models

type RolPermiso struct {
	ID        uint `gorm:"primaryKey;autoIncrement" json:"id"`
	IDRol     uint `gorm:"column:id_rol;not null" json:"id_rol"`
	IDPermiso uint `gorm:"column:id_permiso;not null" json:"id_permiso"`

	Rol     Rol     `gorm:"foreignKey:IDRol;references:ID" json:"-"`
	Permiso Permiso `gorm:"foreignKey:IDPermiso;references:ID" json:"-"`
}

func (RolPermiso) TableName() string {
	return "RolesPermisos"
}
