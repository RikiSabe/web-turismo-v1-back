package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/models"
	"web-turismo-v1/internal/services"
	"web-turismo-v1/internal/types"

	"github.com/gorilla/mux"
)

func ObtenerReservas(w http.ResponseWriter, r *http.Request) {
	var reservas []types.ReservaDTO

	if err := db.GDB.Raw(services.QueryReservasTODO).Scan(&reservas).Error; err != nil {
		http.Error(w, "Error al obtener reservas del usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservas)
}

func ObtenerReserva(w http.ResponseWriter, r *http.Request) {
	id_reserva := mux.Vars(r)["id"]

	query := `
		select 
			r.estado, 
			pt.nombre 
		from "GestReservas"as r 
		join "GestPaquetesTuristicos" as pt on pt.id_paquete_turistico = r.id_paquete
		where r.id_reserva = ?`

	type res struct {
		Nombre string `json:"nombre"`
		Estado bool   `json:"estado"`
	}

	var reserva res
	if err := db.GDB.Raw(query, id_reserva).Scan(&reserva).Error; err != nil {
		http.Error(w, "Error en la consula", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reserva)
}

func ObtenerReservasPorUsuario(w http.ResponseWriter, r *http.Request) {
	id_usuario := mux.Vars(r)["id"]

	type PaqueteInfo struct {
		Nombre      string `json:"nombre"`
		Descripcion string `json:"descripcion"`
	}

	type ReservaDetalle struct {
		ID             uint        `json:"id"`
		Fecha          time.Time   `json:"fecha"`
		Descripcion    string      `json:"descripcion"`
		NumeroPersonas int         `json:"numero_personas"`
		Estado         bool        `json:"estado"`
		Paquete        PaqueteInfo `json:"paquete"`
	}

	var result struct {
		Reservas json.RawMessage `gorm:"column:reservas"`
	}

	if err := db.GDB.Raw(services.QueryReservasByUsuario, id_usuario).Scan(&result).Error; err != nil {
		http.Error(w, "Error al obtener las reservas", http.StatusInternalServerError)
		return
	}

	var reservas []ReservaDetalle
	if err := json.Unmarshal(result.Reservas, &reservas); err != nil {
		http.Error(w, "Error al procesar las reservas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(reservas); err != nil {
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

	b, err := json.MarshalIndent(reserva, "", "  ")
	if err != nil {
		fmt.Println("Error al convertir a JSON:", err)
		return
	}
	fmt.Println(string(b))

	tx := db.GDB.Begin()
	if err := tx.Create(&reserva).Error; err != nil {
		http.Error(w, "Error al crear la reserva", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reserva)
}

func DecisionReserva(w http.ResponseWriter, r *http.Request) {
	id_reserva := mux.Vars(r)["id"]

	type Body struct {
		Estado bool `json:"estado"`
	}

	var decision Body

	if err := json.NewDecoder(r.Body).Decode(&decision); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	var reserva models.Reservas
	if err := db.GDB.First(&reserva, "id_reserva = ?", id_reserva).Error; err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	fmt.Printf("error %v\n", decision.Estado)
	reserva.Estado = decision.Estado

	if err := db.GDB.Save(&reserva).Error; err != nil {
		http.Error(w, "Error al decidir reserva", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reserva)
}
