package controllers

import (
	"encoding/json"
	"net/http"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/services"
	"web-turismo-v1/internal/types"
)

func ObtenerDepartamentos(w http.ResponseWriter, r *http.Request) {
	var departamentos []types.DepartamentoTODO

	err := db.GDB.Raw(services.QueryDepartamentosTODO).Scan(&departamentos).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(departamentos)
}
