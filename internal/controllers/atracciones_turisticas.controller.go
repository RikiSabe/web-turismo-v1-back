package controllers

import (
	"encoding/json"
	"net/http"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/models"
	"web-turismo-v1/internal/services"
	"web-turismo-v1/internal/types"

	"github.com/gorilla/mux"
)

func ObtenerAtraccionesTuristicas(w http.ResponseWriter, r *http.Request) {
	var atraccionesTuristicas []types.AtraccionTuristicaTODO

	err := db.GDB.Raw(services.QueryAtraccionesTuristicasTODO).Scan(&atraccionesTuristicas).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(atraccionesTuristicas); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ObtenerAtraccionTuristica(w http.ResponseWriter, r *http.Request) {
	id_atraccion := mux.Vars(r)["id"]
	var atraccionTuristica types.AtraccionTuristicaTODO

	err := db.GDB.Raw(services.QueryAtraccionesTuristicaUnique, id_atraccion).First(&atraccionTuristica).Error
	if err != nil {
		http.Error(w, "Atraccion Turistica no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&atraccionTuristica); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AgregarAtraccionTuristica(w http.ResponseWriter, r *http.Request) {
	var atraccionTuristica models.AtraccionTuristica

	if err := json.NewDecoder(r.Body).Decode(&atraccionTuristica); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&atraccionTuristica).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Atraccion Turistica", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&atraccionTuristica); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ModificarAtraccionTuristica(w http.ResponseWriter, r *http.Request) {
	id_atraccion := mux.Vars(r)["id"]
	var atraccionExistente models.AtraccionTuristica

	err := db.GDB.Where("id_atracciones = ?", id_atraccion).First(&atraccionExistente).Error
	if err != nil {
		http.Error(w, "Atraccion Turistica no encontrada", http.StatusNotFound)
		return
	}

	var atraccionActualizada types.AtraccionTuristicaTODO
	if err := json.NewDecoder(r.Body).Decode(&atraccionActualizada); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cambios
	atraccionExistente.Tipo = atraccionActualizada.Tipo
	atraccionExistente.Nombre = atraccionActualizada.Nombre
	atraccionExistente.Ubicacion = atraccionActualizada.Ubicacion
	atraccionExistente.Descripcion = atraccionActualizada.Descripcion
	atraccionExistente.Horarios = atraccionActualizada.Horarios
	atraccionExistente.Precio = atraccionActualizada.Precio
	atraccionExistente.Estado = atraccionActualizada.Estado

	if err := db.GDB.Save(&atraccionExistente).Error; err != nil {
		http.Error(w, "Error al actualizar Atraccion Turistica", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&atraccionExistente); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
