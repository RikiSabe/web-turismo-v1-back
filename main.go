package main

import (
	"fmt"
	"log"
	"net/http"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/models"
	"web-turismo-v1/internal/routers"

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
		&models.Rol{},
		&models.Permiso{},
		&models.RolPermiso{},
		&models.Usuario{},
		&models.Agencia{},
		&models.AtraccionTuristica{},

		&models.PaqueteTuristico{},
		&models.PaqueteAtraccion{},
		&models.Reservas{},
	); err != nil {
		log.Fatal("Error al migrar los modelos de la db:", err)
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
