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

func ObtenerUsuarios(w http.ResponseWriter, r *http.Request) {
	var usuarios []types.UsuarioTODO

	err := db.GDB.Raw(services.QueryUsuariosTODO).Scan(&usuarios).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(usuarios); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ObtenerUsuario(w http.ResponseWriter, r *http.Request) {
	id_usuario := mux.Vars(r)["id"]
	var usuario types.UsuarioTODO

	err := db.GDB.Raw(services.QueryUsuarioUnique, id_usuario).First(&usuario).Error
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&usuario); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AgregarUsuario(w http.ResponseWriter, r *http.Request) {
	var usuario models.Usuario

	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&usuario).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Usuario", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&usuario); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ModificarUsuario(w http.ResponseWriter, r *http.Request) {
	id_usuario := mux.Vars(r)["id"]
	var usuarioExistente models.Usuario

	err := db.GDB.Where("id_usuario = ?", id_usuario).First(&usuarioExistente).Error
	if err != nil {
		http.Error(w, "Agencia no encontrada", http.StatusNotFound)
		return
	}

	var usuarioActualizado types.UsuarioTODO
	if err := json.NewDecoder(r.Body).Decode(&usuarioActualizado); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cambios
	usuarioExistente.Rol = usuarioActualizado.Rol
	usuarioExistente.Nombre = usuarioActualizado.Nombre
	usuarioExistente.Apellido = usuarioActualizado.Apellido
	usuarioExistente.Correo = usuarioActualizado.Correo
	usuarioExistente.Telefono = usuarioActualizado.Telefono
	usuarioExistente.Direccion = usuarioActualizado.Direccion
	usuarioExistente.Contra = usuarioActualizado.Contra
	usuarioExistente.Estado = usuarioActualizado.Estado
	usuarioExistente.Foto = usuarioActualizado.Foto
	// Falta eliminar foto del backend

	if err := db.GDB.Save(&usuarioExistente).Error; err != nil {
		http.Error(w, "Error al actualizar usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&usuarioExistente); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
