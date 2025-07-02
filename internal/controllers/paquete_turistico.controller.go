package controllers

import (
	"encoding/json"
	"net/http"
	"time"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/models"
	"web-turismo-v1/internal/services"
	"web-turismo-v1/internal/types"

	"github.com/gorilla/mux"
)

func CrearPaqueteTuristico(w http.ResponseWriter, r *http.Request) {
	var bodyStruct struct {
		Categoria   string    `json:"categoria"`
		Nombre      string    `json:"nombre"`
		Fecha       time.Time `json:"fecha"`
		Descripcion string    `json:"descripcion"`
		Precio      float64   `json:"precio"`
		Duracion    string    `json:"duracion"`
		Salida      string    `json:"salida"`
		Estado      bool      `json:"estado"`
		IDAgencia   uint      `json:"id_agencia"`

		IDAtracciones []uint `json:"id_atracciones"`
	}

	if err := json.NewDecoder(r.Body).Decode(&bodyStruct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var paqueteTuristico = models.PaqueteTuristico{
		Categoria:   bodyStruct.Categoria,
		Nombre:      bodyStruct.Nombre,
		Fecha:       bodyStruct.Fecha,
		Descripcion: bodyStruct.Descripcion,
		Precio:      bodyStruct.Precio,
		Duracion:    bodyStruct.Duracion,
		Salida:      bodyStruct.Salida,
		Estado:      bodyStruct.Estado,
		IDAgencia:   bodyStruct.IDAgencia,
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&paqueteTuristico).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Paquete Turistico", http.StatusInternalServerError)
		return
	}
	for _, idAtraccion := range bodyStruct.IDAtracciones {
		if err := tx.Create(&models.PaqueteAtraccion{IDPaquete: paqueteTuristico.ID, IDAtraccion: idAtraccion}).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al agregar Atracción al Paquete Turístico", http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&bodyStruct); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ObtenerPaquetesTuristicos(w http.ResponseWriter, r *http.Request) {
	var paquetes []types.PaqueteTuristicoTODO

	if err := db.GDB.Raw(services.QueryPaqueteTuristicoTODO).Scan(&paquetes).Error; err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(paquetes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ObtenerPaqueteTuristico(w http.ResponseWriter, r *http.Request) {
	id_paquete := mux.Vars(r)["id"]

	var paquete types.PaqueteTuristicoTODO
	if err := db.GDB.Raw(services.QueryPaqueteTuristicoTODOByID, id_paquete).Scan(&paquete).Error; err != nil {
		http.Error(w, "Error al obtener el paquete turístico", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(paquete); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
