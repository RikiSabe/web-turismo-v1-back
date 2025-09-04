package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/models"
	"web-turismo-v1/internal/routers"
	"web-turismo-v1/internal/seed"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	var err error

	// Cargar el archivo .env
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Error al conectarse a la base de datos: %v", err)
	}

	port := "5000"

	err = db.Connection()
	if err != nil {
		log.Printf("Error al conectar a la base de datos: %v", err)
		return
	}

	if err := db.GDB.AutoMigrate(
		/* migraciones */
		&models.Usuario{},
		&models.Agencia{},
		&models.AtraccionTuristica{},

		&models.PaqueteTuristico{},
		&models.PaqueteAtraccion{},
		&models.Reservas{},

		&models.Departamento{},
		&models.Provincia{},

		&models.FotosAgencia{},
		&models.FotosAtracciones{},

		&models.Categoria{},
		&models.SubCategoria{},

		&models.Dia{},
		&models.Salida{},
		&models.VigenciaDiasConcurrentes{},
	); err != nil {
		log.Fatal("Error al migrar los modelos de la db:", err)
	}

	if err := seed.SeedDatosBolivia(db.GDB, "internal/sql/departamentos_provincias_bolivia.sql"); err != nil {
		log.Printf("Error en el seed de Departamentos: %v", err)
	}

	if err := seed.SeedDatosCategoriaSubCategoria(db.GDB, "internal/sql/categorias_subcategorias.sql"); err != nil {
		log.Printf("Error en el seed de Categorias: %v", err)
	}

	var count int64
	err = db.GDB.Model(&models.Usuario{}).Count(&count).Error
	if err == nil && count == 0 {
		usuario := models.Usuario{
			ID:              1,
			Rol:             "admin",
			Nombre:          "Juan",
			ApellidoPaterno: "Pérez",
			ApellidoMaterno: "Gómez",
			FechaNacimiento: time.Date(1995, 6, 15, 0, 0, 0, 0, time.UTC),
			CI:              "12345678",
			Correo:          "juan.perez@example.com",
			Telefono:        "70000000",
			Contra:          "123456",
			Estado:          true,
			Foto:            "N/A",
			IdUbicacion:     2,
		}
		if err := db.GDB.Create(&usuario).Error; err == nil {
			log.Printf("Primer usuario creador exitosamente: %v", err)
		} else {
			log.Printf("Error al crear el primer usuario")
			return
		}
	}

	r := mux.NewRouter()
	routers.InitEndPoints(r)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	// Iniciar el servidor
	fmt.Printf("Servidor corriendo en puerto: %s\n", port)
	if err := http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk)(r)); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
