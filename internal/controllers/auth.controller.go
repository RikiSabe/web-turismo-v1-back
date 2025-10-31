package controllers

import (
	"encoding/json"
	"net/http"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/models"
	"web-turismo-v1/internal/types"
)

type auth struct{}

var Auth auth

func (auth) AuthLoginWeb(w http.ResponseWriter, r *http.Request) {
	var usuario types.Usuario
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var usuarioExistente models.Usuario
	if err := db.GDB.Where("correo = ?", usuario.Correo).First(&usuarioExistente).Error; err != nil {
		http.Error(w, "Credenciales incorrectas", http.StatusInternalServerError)
		return
	}

	if !usuarioExistente.Estado {
		http.Error(w, "Usuario no habilitado", http.StatusInternalServerError)
		return
	}

	if usuarioExistente.Contra != usuario.Contra {
		http.Error(w, "Credenciales incorrectas", http.StatusInternalServerError)
		return
	}

	respuesta := types.RespuestaUsuario{
		ID:  usuarioExistente.ID,
		Rol: usuarioExistente.Rol,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&respuesta); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (auth) AuthRegisterWeb(w http.ResponseWriter, r *http.Request) {
	var nuevoUsuario types.NewUsuario

	if err := json.NewDecoder(r.Body).Decode(&nuevoUsuario); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	usuario := models.Usuario{
		Nombre:          nuevoUsuario.Nombre,
		ApellidoPaterno: nuevoUsuario.ApellidoPaterno,
		ApellidoMaterno: nuevoUsuario.ApellidoMaterno,
		Correo:          nuevoUsuario.Correo,
		Contra:          nuevoUsuario.Contra,
		IdUbicacion:     1,
		Rol:             "turista",
		Estado:          true,
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&usuario).Error; err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx.Commit()

	respuesta := types.RespuestaUsuario{
		ID:  usuario.ID,
		Rol: usuario.Rol,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&respuesta); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
