package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/models"
	"web-turismo-v1/internal/services"
	"web-turismo-v1/internal/types"

	"github.com/gorilla/mux"
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
	id_agencia := mux.Vars(r)["id"]
	var agencia types.AgenciaTODO

	err := db.GDB.Raw(services.QueryAgenciaUnique, id_agencia).Scan(&agencia).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agencia)
}

func AgregarAgencia(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error al parsear el formulario", http.StatusInternalServerError)
		return
	}

	idEncargado := uint(0)
	if idEncargadoStr := r.FormValue("id_encargado"); idEncargadoStr != "" {
		parsed, err := strconv.ParseUint(idEncargadoStr, 10, 64)
		if err != nil {
			http.Error(w, "id_encargado debe ser un numero valido", http.StatusBadRequest)
			return
		}
		idEncargado = uint(parsed)
	}

	idUbicacion := uint(0)
	if idUbicacionStr := r.FormValue("id_departamento"); idUbicacionStr != "" {
		parsed, err := strconv.ParseUint(idUbicacionStr, 10, 64)
		if err != nil {
			http.Error(w, "id_ubicacion debe ser un número válido", http.StatusBadRequest)
			return
		}
		idUbicacion = uint(parsed)
	}

	nuevaAgencia := models.Agencia{
		Nombre:         r.FormValue("nombre"),
		Telefono:       r.FormValue("telefono"),
		Correo:         r.FormValue("correo"),
		IdEncargado:    idEncargado,
		Descripcion:    r.FormValue("descripcion"),
		IdDepartamento: idUbicacion,
		Direccion:      r.FormValue("direccion"),
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&nuevaAgencia).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al guardar la agencia", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nuevaAgencia)
}

func ModificarAgencia(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	var agencia models.Agencia
	if err := db.GDB.Where("id_agencia = ?", idStr).First(&agencia).Error; err != nil {
		http.Error(w, "Agencia no encontrada", http.StatusNotFound)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Error al parsear el formulario: "+err.Error(), http.StatusBadRequest)
		return
	}

	idEncargado := uint(0)
	if idEncargadoStr := r.FormValue("id_encargado"); idEncargadoStr != "" {
		parsed, err := strconv.ParseUint(idEncargadoStr, 10, 64)
		if err != nil {
			http.Error(w, "id_encargado debe ser un numero valido", http.StatusBadRequest)
			return
		}
		idEncargado = uint(parsed)
	}

	idUbicacion := uint(0)
	if idUbicacionStr := r.FormValue("id_ubicacion"); idUbicacionStr != "" {
		parsed, err := strconv.ParseUint(idUbicacionStr, 10, 64)
		if err != nil {
			http.Error(w, "id_ubicacion debe ser un número válido", http.StatusBadRequest)
			return
		}
		idUbicacion = uint(parsed)
	}

	agencia.Nombre = r.FormValue("nombre")
	agencia.Direccion = r.FormValue("direccion")
	agencia.Telefono = r.FormValue("telefono")
	agencia.Correo = r.FormValue("correo")
	agencia.Estado = r.FormValue("estado") == "true"
	agencia.IdEncargado = idEncargado
	agencia.Descripcion = r.FormValue("descripcion")
	agencia.IdDepartamento = idUbicacion

	if err := db.GDB.Save(&agencia).Error; err != nil {
		http.Error(w, "Error al actualizar la agencia", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agencia)
}
