package controllers

import (
	"encoding/json"
	"net/http"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/services"
	"web-turismo-v1/internal/types"
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
	// code
}

func AgregarAgencia(w http.ResponseWriter, r *http.Request) {
	// code
}

func ModificarAgencia(w http.ResponseWriter, r *http.Request) {
	// code
}
