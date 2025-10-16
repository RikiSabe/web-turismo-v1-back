package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/models"

	"github.com/gorilla/mux"
)

func AsignarQR(w http.ResponseWriter, r *http.Request) {
	id_encargado := mux.Vars(r)["id"]
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error al parsear el formulario", http.StatusInternalServerError)
		return
	}

	var encargadoQR models.EncargadoQR
	db.GDB.Where("id_usuario = ?", id_encargado).First(&encargadoQR)

	if encargadoQR.Foto != "" && encargadoQR.Foto != "N/A" {
		if _, err := os.Stat(encargadoQR.Foto); err == nil {
			_ = os.Remove(encargadoQR.Foto)
		}
	}

	direccionFoto := "N/A"

	outputDir := "internal/images/qr"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		http.Error(w, "Error al crear la carpeta de imágenes", http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("QR")
	if err == nil {
		defer file.Close()

		nombreFoto := fmt.Sprintf("encargado_qr_%s%s", id_encargado, filepath.Ext(handler.Filename))
		rutaFoto := "internal/images/qr/" + nombreFoto

		if err := os.MkdirAll("internal/images/qr", os.ModePerm); err != nil {
			http.Error(w, "Error al crear carpeta de QR", http.StatusInternalServerError)
			return
		}

		outFile, err := os.Create(rutaFoto)
		if err != nil {
			http.Error(w, "Error al guardar la foto", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(w, "Error al escribir la foto", http.StatusInternalServerError)
			return
		}

		direccionFoto = rutaFoto
	} else if err != http.ErrMissingFile {
		http.Error(w, "Error al obtener la foto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	encargadoQR.Foto = direccionFoto
	encargadoQR.IdUsuario = parseUint(id_encargado)

	tx := db.GDB.Begin()
	if encargadoQR.ID != 0 {
		if err := tx.Save(&encargadoQR).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al actualizar QR", http.StatusInternalServerError)
			return
		}
	} else {
		if err := tx.Create(&encargadoQR).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al guardar QR", http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(encargadoQR)
}

func ObtenerQR(w http.ResponseWriter, r *http.Request) {
	id_encargado := mux.Vars(r)["id"]

	var encargadoQR models.EncargadoQR
	if err := db.GDB.Where("id_usuario = ?", id_encargado).First(&encargadoQR).Error; err != nil {
		http.Error(w, "QR no encontrado", http.StatusNotFound)
		return
	}

	if encargadoQR.Foto == "" || encargadoQR.Foto == "N/A" {
		http.Error(w, "Este encargado no tiene QR asignado", http.StatusNotFound)
		return
	}

	base64Str, err := encodeImageToBase64(encargadoQR.Foto)
	if err != nil {
		http.Error(w, "Error al codificar imagen", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"id_usuario": id_encargado,
		"foto":       base64Str,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func parseUint(s string) uint {
	var id uint
	fmt.Sscanf(s, "%d", &id)
	return id
}
