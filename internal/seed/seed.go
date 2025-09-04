package seed

import (
	"fmt"
	"os"
	"web-turismo-v1/internal/models"

	"gorm.io/gorm"
)

func SeedDatosBolivia(db *gorm.DB, sqlFilePath string) error {
	var countDepartamentos int64
	if err := db.Model(&models.Departamento{}).Count(&countDepartamentos).Error; err != nil {
		return fmt.Errorf("error al contar departamentos: %v", err)
	}
	if countDepartamentos > 0 {
		fmt.Println("Departamentos ya precargados, omitiendo seed")
		return nil
	}

	// Leer el archivo .sql
	sqlBytes, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return fmt.Errorf("error al leer archivo SQL: %v", err)
	}

	sqlText := string(sqlBytes)

	// Ejecutar el contenido del archivo SQL
	if err := db.Exec(sqlText).Error; err != nil {
		return fmt.Errorf("error al ejecutar SQL: %v", err)
	}

	fmt.Println("Seed de Departamentos ejecutado correctamente desde archivo SQL")
	return nil
}

func SeedDatosCategoriaSubCategoria(db *gorm.DB, sqlFilePath string) error {
	var countCategorias int64
	if err := db.Model(&models.Categoria{}).Count(&countCategorias).Error; err != nil {
		return fmt.Errorf("error al contar categorias: %v", err)
	}
	if countCategorias > 0 {
		fmt.Println("Categorias ya precargadas, omitiendo seed")
		return nil
	}

	// Leer el archivo .sql
	sqlBytes, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return fmt.Errorf("error al leer archivo SQL: %v", err)
	}

	sqlText := string(sqlBytes)

	// Ejecutar el contenido del archivo SQL
	if err := db.Exec(sqlText).Error; err != nil {
		return fmt.Errorf("error al ejecutar SQL: %v", err)
	}

	fmt.Println("Seed de Categoria ejecutado correctamente desde archivo SQL")
	return nil
}

func SeedDiaSemana(db *gorm.DB, sqlFilePath string) error {
	var countDias int64
	if err := db.Model(&models.Dia{}).Count(&countDias).Error; err != nil {
		return fmt.Errorf("error al contar dias: %v", err)
	}
	if countDias > 0 {
		fmt.Println("Dias ya precargados, omitiendo seed")
		return nil
	}

	// Leer el archivo .sql
	sqlBytes, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return fmt.Errorf("error al leer archivo SQL: %v", err)
	}

	sqlText := string(sqlBytes)

	// Ejecutar el contenido del archivo SQL
	if err := db.Exec(sqlText).Error; err != nil {
		return fmt.Errorf("error al ejecutar SQL: %v", err)
	}

	fmt.Println("Seed de Dias ejecutado correctamente desde archivo SQL")
	return nil
}
