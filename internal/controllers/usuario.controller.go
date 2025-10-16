package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/models"
	"web-turismo-v1/internal/services"
	"web-turismo-v1/internal/types"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Function added
func encodeImageToBase64(path string) (string, error) {
	imgBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	mimeType := "image/jpeg"
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png":
		mimeType = "image/png"
	case ".jpg":
		mimeType = "image/jpg"
	case ".jpeg":
		mimeType = "image/jpeg"
	default:
		mimeType = "application/octet-stream"
	}

	encoded := base64.StdEncoding.EncodeToString(imgBytes)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}

func ObtenerUsuarios(w http.ResponseWriter, r *http.Request) {
	var usuarios []types.UsuarioResumen

	err := db.GDB.Raw(services.QueryUsuariosResumen).Scan(&usuarios).Error
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

func ObtenerUsuarioMenu(w http.ResponseWriter, r *http.Request) {
	id_usuario := mux.Vars(r)["id"]
	var usuario types.UsuarioMenu

	err := db.GDB.Raw(services.QueryUsuarioMenu, id_usuario).First(&usuario).Error
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	fotoBase64 := ""
	if usuario.Foto != "N/A" {
		if encoded, err := encodeImageToBase64(usuario.Foto); err == nil {
			fotoBase64 = encoded
		}
	}

	if fotoBase64 == "" {
		fotoBase64 = "N/A"
	}

	usuario.Foto = fotoBase64

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&usuario); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ObtenerUsuario(w http.ResponseWriter, r *http.Request) {
	id_usuario := mux.Vars(r)["id"]
	var usuario types.UsuarioUnique

	err := db.GDB.Raw(services.QueryUsuarioUnique, id_usuario).First(&usuario).Error
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	fotoBase64 := ""
	if usuario.Foto != "N/A" {
		if encoded, err := encodeImageToBase64(usuario.Foto); err == nil {
			fotoBase64 = encoded
		}
	}
	if fotoBase64 == "" {
		fotoBase64 = "N/A"
	}

	usuario.Foto = fotoBase64

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&usuario); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AgregarUsuario(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error al parsear el formulario", http.StatusInternalServerError)
		return
	}

	direccionFoto := "N/A"

	outputDir := "internal/images/usuarios"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		http.Error(w, "Error al crear la carpeta de imágenes", http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("foto")
	if err == nil {
		defer file.Close()

		var count int64
		db.GDB.Model(&models.Usuario{}).Count(&count)
		nombreFoto := fmt.Sprintf("foto_usuario_%d%s", count+1, filepath.Ext(handler.Filename))

		rutaFoto := "internal/images/usuarios/" + nombreFoto
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

	fechaNacimientoStr := r.FormValue("fecha_nacimiento")
	fechaNacimiento, err := time.Parse("2006-01-02", fechaNacimientoStr)
	if err != nil {
		http.Error(w, "Fecha de nacimiento inválida. Formato esperado: YYYY-MM-DD", http.StatusBadRequest)
		return
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

	nuevoUsuario := models.Usuario{
		Rol:             r.FormValue("rol"),
		Nombre:          r.FormValue("nombre"),
		ApellidoPaterno: r.FormValue("apellido_paterno"),
		ApellidoMaterno: r.FormValue("apellido_materno"),
		FechaNacimiento: fechaNacimiento,
		CI:              r.FormValue("ci"),
		Correo:          r.FormValue("correo"),
		Telefono:        r.FormValue("telefono"),
		Contra:          r.FormValue("contra"),
		IdUbicacion:     idUbicacion,
		Estado:          true,
		Foto:            direccionFoto,
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&nuevoUsuario).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al guardar el usuario", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nuevoUsuario)
}

func ModificarUsuario(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	var usuario models.Usuario
	if err := db.GDB.Where("id_usuario = ?", idStr).First(&usuario).Error; err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Error al parsear el formulario: "+err.Error(), http.StatusBadRequest)
		return
	}

	fechaNacimientoStr := r.FormValue("fecha_nacimiento")
	fechaNacimiento, err := time.Parse("2006-01-02", fechaNacimientoStr)
	if err != nil {
		http.Error(w, "Formato de fecha inválido. Debe ser YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	usuario.Rol = r.FormValue("rol")
	usuario.Nombre = r.FormValue("nombre")
	usuario.ApellidoPaterno = r.FormValue("apellido_paterno")
	usuario.ApellidoMaterno = r.FormValue("apellido_materno")
	usuario.FechaNacimiento = fechaNacimiento
	usuario.CI = r.FormValue("ci")
	usuario.Correo = r.FormValue("correo")
	usuario.Telefono = r.FormValue("telefono")
	usuario.Estado = r.FormValue("estado") == "true"

	file, handler, err := r.FormFile("foto")
	if err == nil {
		defer file.Close()

		if usuario.Foto != "" && usuario.Foto != "N/A" {
			_ = os.Remove(usuario.Foto) // ignorar error si no existe
		}

		var count int64
		db.GDB.Model(&models.Usuario{}).Count(&count)
		nombreFoto := fmt.Sprintf("foto_usuario_%d%s", count+1, filepath.Ext(handler.Filename))
		rutaFoto := "internal/images/usuarios/" + nombreFoto

		outFile, err := os.Create(rutaFoto)
		if err != nil {
			http.Error(w, "Error al guardar la nueva foto", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, file); err != nil {
			http.Error(w, "Error al escribir la nueva foto", http.StatusInternalServerError)
			return
		}

		usuario.Foto = rutaFoto
	} else if err != http.ErrMissingFile {
		http.Error(w, "Error al procesar la foto: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.GDB.Save(&usuario).Error; err != nil {
		http.Error(w, "Error al actualizar el usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// Steps
func ObtenerUsuarioDatosPersonales(w http.ResponseWriter, r *http.Request) {
	id_usuario := mux.Vars(r)["id"]

	var datosPersonales types.UsuarioDatosPersonales

	query := `
			SELECT 
				u.nombre,
				u.apellido_paterno,
				u.apellido_materno,
				u.telefono,
				u.fecha_nacimiento,
				u.id_ubicacion,
				p.id_departamento
			FROM "GestUsuarios" u
			LEFT JOIN provincias p ON u.id_ubicacion = p.id_provincia
			WHERE u.id_usuario = $1 AND u.estado = true`

	if err := db.GDB.Raw(query, id_usuario).Scan(&datosPersonales).Error; err != nil {
		http.Error(w, "Error al obtener datos personales del usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(datosPersonales)
}

func ObtenerUsuarioDatosPrivados(w http.ResponseWriter, r *http.Request) {
	id_usuario := mux.Vars(r)["id"]
	var datosPrivados types.UsuarioDatosPrivados

	query := `
		SELECT 
			rol,
			ci,
			correo,
			contra,
			estado
		FROM "GestUsuarios"
		WHERE id_usuario = $1`

	if err := db.GDB.Raw(query, id_usuario).Scan(&datosPrivados).Error; err != nil {
		http.Error(w, "Error al obtener datos privados del usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(datosPrivados)
}

func ObtenerUsuarioFoto(w http.ResponseWriter, r *http.Request) {
	id_usuario := mux.Vars(r)["id"]
	var usuarioFoto types.UsuarioFoto

	query := `
			SELECT foto
			FROM "GestUsuarios"
			WHERE id_usuario = $1`

	if err := db.GDB.Raw(query, id_usuario).Scan(&usuarioFoto).Error; err != nil {
		http.Error(w, "Error al obtener foto del usuario", http.StatusInternalServerError)
		return
	}

	fotoBase64 := ""
	if usuarioFoto.Foto != "N/A" {
		if encoded, err := encodeImageToBase64(usuarioFoto.Foto); err == nil {
			fotoBase64 = encoded
		}
	}
	if fotoBase64 == "" {
		fotoBase64 = "N/A"
	}
	usuarioFoto.Foto = fotoBase64

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarioFoto)
}

// Mod Steps
func ModificarUsuarioDatosPersonales(w http.ResponseWriter, r *http.Request) {
	id_usuario := mux.Vars(r)["id"]

	var usuarioExistente models.Usuario
	err := db.GDB.
		Where("id_usuario = ?", id_usuario).
		First(&usuarioExistente).
		Error
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	var nuevosDatos types.UsuarioDatosPersonales
	if err := json.NewDecoder(r.Body).Decode(&nuevosDatos); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usuarioExistente.Nombre = nuevosDatos.Nombre
	usuarioExistente.ApellidoPaterno = nuevosDatos.ApellidoPaterno
	usuarioExistente.ApellidoMaterno = nuevosDatos.ApellidoMaterno
	usuarioExistente.Telefono = nuevosDatos.Telefono
	usuarioExistente.FechaNacimiento = nuevosDatos.FechaNacimiento
	usuarioExistente.IdUbicacion = nuevosDatos.IdUbicacion

	if err := db.GDB.Save(&usuarioExistente).Error; err != nil {
		http.Error(w, "Error al modificar datos personales del usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarioExistente)
}

func ModificarUsuarioDatosPrivados(w http.ResponseWriter, r *http.Request) {
	id_usuario := mux.Vars(r)["id"]

	var usuarioExistente models.Usuario
	err := db.GDB.
		Where("id_usuario = ?", id_usuario).
		First(&usuarioExistente).
		Error
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	var nuevosDatos types.UsuarioDatosPrivados
	if err := json.NewDecoder(r.Body).Decode(&nuevosDatos); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usuarioExistente.Rol = nuevosDatos.Rol
	usuarioExistente.CI = nuevosDatos.CI
	usuarioExistente.Correo = nuevosDatos.Correo
	usuarioExistente.Contra = nuevosDatos.Contra
	usuarioExistente.Estado = nuevosDatos.Estado

	if err := db.GDB.Save(&usuarioExistente).Error; err != nil {
		http.Error(w, "Error al modificar datos privados del usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarioExistente)
}

func ModificarUsuarioFoto(w http.ResponseWriter, r *http.Request) {
	id_usuario := mux.Vars(r)["id"]

	var usuarioExistente models.Usuario
	err := db.GDB.Where("id_usuario = ?", id_usuario).First(&usuarioExistente).Error
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error al parsear el formulario: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("foto")
	if err != nil {
		if err == http.ErrMissingFile {
			http.Error(w, "No se recibió ninguna foto", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al obtener la foto: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer file.Close()

	if usuarioExistente.Foto != "" && usuarioExistente.Foto != "N/A" {
		_ = os.Remove(usuarioExistente.Foto) // Ignorar error si no existe
	}

	nombreFoto := fmt.Sprintf("foto_usuario_%s%s", uuid.New().String(), filepath.Ext(handler.Filename))
	rutaFoto := "internal/images/usuarios/" + nombreFoto

	outFile, err := os.Create(rutaFoto)
	if err != nil {
		http.Error(w, "Error al guardar la foto en disco", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		http.Error(w, "Error al escribir la foto", http.StatusInternalServerError)
		return
	}

	usuarioExistente.Foto = rutaFoto
	if err := db.GDB.Save(&usuarioExistente).Error; err != nil {
		http.Error(w, "Error al actualizar la foto del usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Foto actualizada exitosamente", "foto": rutaFoto})
}
