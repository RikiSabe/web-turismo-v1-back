package controllers

import (
	"encoding/json"
	"net/http"
	"web-turismo-v1/internal/db"

	"github.com/gorilla/mux"
)

func ObtenerDetallesComprobante(w http.ResponseWriter, r *http.Request) {
	id_usuario := mux.Vars(r)["id_usuario"]
	id_reserva := mux.Vars(r)["id_reserva"]

	type Response = struct {
		Foto string `json:"foto"`
	}
	var fotoComprobante Response
	err := db.GDB.
		Table("comprobantes").
		Select("foto").
		Where("id_reserva = ? AND id_usuario = ?", id_reserva, id_usuario).
		Scan(&fotoComprobante).Error

	if err != nil {
		http.Error(w, "Error al obtener el comprobante", http.StatusInternalServerError)
		return
	}

	fotoBase64 := ""
	if fotoComprobante.Foto != "N/A" {
		if encoded, err := encodeImageToBase64(fotoComprobante.Foto); err == nil {
			fotoBase64 = encoded
		}
	}
	if fotoBase64 == "" {
		fotoBase64 = "N/A"
	}

	fotoComprobante.Foto = fotoBase64

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&fotoComprobante)
}
