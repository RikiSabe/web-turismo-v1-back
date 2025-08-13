package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/models"
	"web-turismo-v1/internal/services"
	"web-turismo-v1/internal/types"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/datatypes"
)

type Foto struct {
	ID    uint   `json:"id"`
	Foto  string `json:"foto"`
	Orden uint   `json:"orden"`
}

func ObtenerAtraccionesTuristicas(w http.ResponseWriter, r *http.Request) {
	var atraccionesTuristicas []types.AtraccionTuristicaTODO

	err := db.GDB.Raw(services.QueryAtraccionesTuristicasTODO).Scan(&atraccionesTuristicas).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(atraccionesTuristicas); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ObtenerAtraccionTuristica(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var atraccion types.AtraccionTuristicaUnique

	err := db.GDB.Raw(services.QueryAtraccionesTuristicaUnique, id).Scan(&atraccion).Error
	if err != nil {
		http.Error(w, "Atracción no encontrada", http.StatusNotFound)
		return
	}

	fotos := atraccion.Fotos.Data()

	for i, foto := range fotos {
		base64img, err := encodeImageToBase64(foto.Foto)
		if err != nil {
			fmt.Printf("Error codificando una foto")

		}
		fotos[i].Foto = base64img
	}
	atraccion.Fotos = datatypes.NewJSONType(fotos)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&atraccion)
}

func AgregarAtraccionTuristica(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(50 << 20) // 50 MB
	if err != nil {
		http.Error(w, "Error al parsear el formulario", http.StatusInternalServerError)
		return
	}

	idEncargado, err := strconv.ParseUint(r.FormValue("id_encargado"), 10, 64)
	if err != nil {
		http.Error(w, "id_encargado inválido", http.StatusBadRequest)
		return
	}

	idUbicacion, err := strconv.ParseUint(r.FormValue("id_ubicacion"), 10, 64)
	if err != nil {
		http.Error(w, "id_ubicacion inválido", http.StatusBadRequest)
		return
	}

	precio, err := strconv.ParseFloat(r.FormValue("precio"), 64)
	if err != nil {
		http.Error(w, "precio inválido", http.StatusBadRequest)
		return
	}

	nuevaAtraccion := models.AtraccionTuristica{
		IdEncargado:     uint(idEncargado),
		IdUbicacion:     uint(idUbicacion),
		Categoria:       r.FormValue("categoria"),
		Nombre:          r.FormValue("nombre"),
		Direccion:       r.FormValue("direccion"),
		Descripcion:     r.FormValue("descripcion"),
		HorarioApertura: r.FormValue("horario_apertura"),
		HorarioCierre:   r.FormValue("horario_cierre"),
		Precio:          precio,
		Estado:          true,
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&nuevaAtraccion).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al guardar la atracción turística", http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["fotos[]"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al abrir una foto: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		ext := filepath.Ext(fileHeader.Filename)
		fmt.Printf(fileHeader.Filename)
		OrdenFoto, err := strconv.Atoi(strings.TrimSuffix((fileHeader.Filename), ext))
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al procesar el orden de la foto", http.StatusBadRequest)
			return
		}
		// Generar nombre único
		nombreFoto := fmt.Sprintf("foto_atraccion_%s%s", uuid.New().String(), filepath.Ext(fileHeader.Filename))
		rutaFoto := "internal/images/atracciones/" + nombreFoto

		outFile, err := os.Create(rutaFoto)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al guardar la foto en disco", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, file); err != nil {
			tx.Rollback()
			http.Error(w, "Error al escribir la foto", http.StatusInternalServerError)
			return
		}

		foto := models.FotosAtracciones{
			IdAtraccion: nuevaAtraccion.ID,
			Foto:        rutaFoto,
			Orden:       uint(OrdenFoto),
		}

		if err := tx.Create(&foto).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al guardar la foto en base de datos", http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nuevaAtraccion)
}

func ModificarAtraccionTuristica(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	var atraccion models.AtraccionTuristica
	if err := db.GDB.Where("id_atraccion = ?", idStr).First(&atraccion).Error; err != nil {
		http.Error(w, "Atraccion no encontrada", http.StatusNotFound)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Error al parsear el formulario: "+err.Error(), http.StatusBadRequest)
		return
	}

	nuevoPrecio := float64(0)

	precioStr := r.FormValue("precio")
	if precioStr != "" {
		parsed, err := strconv.ParseFloat(precioStr, 64)
		if err != nil {
			http.Error(w, "El precio debe ser un número válido", http.StatusBadRequest)
			return
		}
		nuevoPrecio = parsed
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

	idEncargado := uint(0)
	if idEncargadoStr := r.FormValue("id_encargado"); idEncargadoStr != "" {
		parsed, err := strconv.ParseUint(idEncargadoStr, 10, 64)
		if err != nil {
			http.Error(w, "id_encargado debe ser un número válido", http.StatusBadRequest)
			return
		}
		idEncargado = uint(parsed)
	}

	atraccion.Nombre = r.FormValue("nombre")
	atraccion.Categoria = r.FormValue("tipo")
	atraccion.Direccion = r.FormValue("ubicacion")
	atraccion.Descripcion = r.FormValue("descripcion")
	atraccion.HorarioApertura = r.FormValue("horario_apertura")
	atraccion.HorarioCierre = r.FormValue("horario_cierre")
	atraccion.Precio = float64(nuevoPrecio)
	atraccion.Estado = r.FormValue("estado") == "true"
	atraccion.IdEncargado = idEncargado
	atraccion.IdUbicacion = idUbicacion

	if err := db.GDB.Save(&atraccion).Error; err != nil {
		http.Error(w, "Error al actualizar el usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(atraccion)
}

func ObtenerEncargadoAtraccionTuristica(w http.ResponseWriter, r *http.Request) {
	id_encargado := mux.Vars(r)["id"]
	var usuario types.UsuarioEncargado

	err := db.GDB.Raw(services.QueryEncargadoAtraccionTuristicas, id_encargado).Scan(&usuario).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	fotoBase64 := ""
	if usuario.SRC != "N/A" {
		if encoded, err := encodeImageToBase64(usuario.SRC); err == nil {
			fotoBase64 = encoded
		}
	}

	if fotoBase64 == "" {
		fotoBase64 = "N/A"
	}

	usuario.SRC = fotoBase64

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&usuario)
}
