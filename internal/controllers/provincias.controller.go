package controllers

import (
	"encoding/json"
	"net/http"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/services"
	"web-turismo-v1/internal/types"

	"github.com/gorilla/mux"
)

func ObtenerProvinciasbyDepartamento(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var provincias []types.ProvinciasTODO

	err := db.GDB.Raw(services.QueryProvinciasTODO, id).Scan(&provincias).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(provincias); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
