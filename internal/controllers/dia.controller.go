package controllers

import (
	"encoding/json"
	"net/http"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/models"
)

func ObtenerDias(w http.ResponseWriter, r *http.Request) {
	var dias []models.Dia

	err := db.GDB.
		Table("dias").
		Select("id_dia, nombre").
		Scan(&dias).
		Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dias)
}
