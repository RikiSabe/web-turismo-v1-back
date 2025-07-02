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

func ObtenerReservasUsuario(w http.ResponseWriter, r *http.Request) {
	idUsuario := mux.Vars(r)["id"]
	var Reservas []types.ReservaTODO

	err := db.GDB.Raw(services.QueryReservasUsuarioTODO, idUsuario).Scan(&Reservas).Error
	if err != nil {
		http.Error(w, "Error al obtener reservas del usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Reservas); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HacerReserva(w http.ResponseWriter, r *http.Request) {
	var reserva models.Reservas
	if err := json.NewDecoder(r.Body).Decode(&reserva); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	if err := db.GDB.Create(&reserva).Error; err != nil {
		http.Error(w, "Error al crear la reserva", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(reserva); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
