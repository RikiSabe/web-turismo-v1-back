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

func ObtenerRoles(w http.ResponseWriter, r *http.Request) {
	var roles []types.RolTODO

	err := db.GDB.Raw(services.QueryRolesTODO).Scan(&roles).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(roles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ObtenerRol(w http.ResponseWriter, r *http.Request) {
	id_rol := mux.Vars(r)["id"]
	var rol types.RolTODO

	err := db.GDB.Raw(services.QueryRolUnique, id_rol).First(&rol).Error
	if err != nil {
		http.Error(w, "Rol no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&rol); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AgregarRol(w http.ResponseWriter, r *http.Request) {
	var rol models.Rol

	if err := json.NewDecoder(r.Body).Decode(&rol); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&rol).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Rol", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&rol); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ModificarRol(w http.ResponseWriter, r *http.Request) {
	id_rol := mux.Vars(r)["id"]
	var rolExistente models.Rol

	err := db.GDB.Where("id_rol = ?", id_rol).First(&rolExistente).Error
	if err != nil {
		http.Error(w, "Agencia no encontrada", http.StatusNotFound)
		return
	}

	var rolActualizado types.RolTODO
	if err := json.NewDecoder(r.Body).Decode(&rolActualizado); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cambios
	rolExistente.Nombre = rolActualizado.Nombre
	rolExistente.Estado = rolActualizado.Estado

	if err := db.GDB.Save(&rolExistente).Error; err != nil {
		http.Error(w, "Error al actualizar Rol", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&rolExistente); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
