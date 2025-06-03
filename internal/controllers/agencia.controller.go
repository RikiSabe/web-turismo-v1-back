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

func ObtenerAgencias(w http.ResponseWriter, r *http.Request) {
	var agencias []types.AgenciaTODO

	err := db.GDB.Raw(services.QueryAgenciasTODO).Scan(&agencias).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(agencias); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ObtenerAgencia(w http.ResponseWriter, r *http.Request) {
	id_agencia := mux.Vars(r)["id"]
	var agencia types.AgenciaTODO

	err := db.GDB.Raw(services.QueryAgenciaUnique, id_agencia).First(&agencia).Error
	if err != nil {
		http.Error(w, "Agencia no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&agencia); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AgregarAgencia(w http.ResponseWriter, r *http.Request) {
	var agencia models.Agencia

	if err := json.NewDecoder(r.Body).Decode(&agencia); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&agencia).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Agencia", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&agencia); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ModificarAgencia(w http.ResponseWriter, r *http.Request) {
	id_agencia := mux.Vars(r)["id"]
	var agenciaExistente models.Agencia

	err := db.GDB.Where("id_agencia = ?", id_agencia).First(&agenciaExistente).Error
	if err != nil {
		http.Error(w, "Agencia no encontrada", http.StatusNotFound)
		return
	}

	var agenciaActualizada types.AgenciaTODO
	if err := json.NewDecoder(r.Body).Decode(&agenciaActualizada); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cambios
	agenciaExistente.Nombre = agenciaActualizada.Nombre
	agenciaExistente.Direccion = agenciaActualizada.Direccion
	agenciaExistente.Telefono = agenciaActualizada.Telefono
	agenciaExistente.CorreoElectronico = agenciaActualizada.CorreoElectronico
	agenciaExistente.Estado = agenciaActualizada.Estado

	if err := db.GDB.Save(&agenciaExistente).Error; err != nil {
		http.Error(w, "Error al actualizar agencia", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&agenciaExistente); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
