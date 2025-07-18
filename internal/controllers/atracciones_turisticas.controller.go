package controllers

import (
	"encoding/json"
	"net/http"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/services"
	"web-turismo-v1/internal/types"
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
	// code
}

func AgregarAtraccionTuristica(w http.ResponseWriter, r *http.Request) {
	// code
}

func ModificarAtraccionTuristica(w http.ResponseWriter, r *http.Request) {
	// code
}
