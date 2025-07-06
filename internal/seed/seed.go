package seed

import (
	"fmt"
	"io/ioutil"
	"web-turismo-v1/internal/models"

	"gorm.io/gorm"
)

func SeedDatosBolivia(db *gorm.DB, sqlFilePath string) error {
	var count int64
	if err := db.Model(&models.Departamento{}).Count(&count).Error; err != nil {
		return fmt.Errorf("error al contar departamentos: %v", err)
	}
	if count > 0 {
		fmt.Println("Departamentos ya precargados, omitiendo seed")
		return nil
	}

	// Leer el archivo .sql
	sqlBytes, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		return fmt.Errorf("error al leer archivo SQL: %v", err)
	}

	sqlText := string(sqlBytes)

	// Ejecutar el contenido del archivo SQL
	if err := db.Exec(sqlText).Error; err != nil {
		return fmt.Errorf("error al ejecutar SQL: %v", err)
	}

	fmt.Println("Seed ejecutado correctamente desde archivo SQL")
	return nil
}
