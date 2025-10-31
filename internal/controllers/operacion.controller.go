package controllers

import (
	"encoding/json"
	"net/http"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/models"
)

func HacerOperacion(w http.ResponseWriter, r *http.Request) {
	var operacion models.Operacion

	if err := json.NewDecoder(r.Body).Decode(&operacion); err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&operacion).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al crear la operación", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(operacion)
}
