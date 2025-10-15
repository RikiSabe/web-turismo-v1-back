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
	var agencia types.AgenciaUnique

	err := db.GDB.Raw(services.QueryAgenciaUnique, id_agencia).Scan(&agencia).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	fotos := agencia.Fotos.Data()

	if len(fotos) == 0 {
		agencia.Fotos = datatypes.NewJSONType([]types.Foto{})
	} else {
		for i, foto := range fotos {
			base64img, err := encodeImageToBase64(foto.Foto)
			if err != nil {
				fmt.Printf("Error codificando una foto")

			}
			fotos[i].Foto = base64img
		}
		agencia.Fotos = datatypes.NewJSONType(fotos)
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
		Direccion:      r.FormValue("direccion"),
		Telefono:       r.FormValue("telefono"),
		Correo:         r.FormValue("correo"),
		Estado:         r.FormValue("estado") == "true",
		IdEncargado:    idEncargado,
		Descripcion:    r.FormValue("descripcion"),
		IdDepartamento: idUbicacion,
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&nuevaAgencia).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al guardar la agencia", http.StatusInternalServerError)
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
		nombreFoto := fmt.Sprintf("foto_agencia_%s%s", uuid.New().String(), filepath.Ext(fileHeader.Filename))
		rutaFoto := "internal/images/agencias/" + nombreFoto

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

		foto := models.FotosAgencia{
			IdAgencia: nuevaAgencia.ID,
			Foto:      rutaFoto,
			Orden:     uint(OrdenFoto),
		}

		if err := tx.Create(&foto).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al guardar la foto en base de datos", http.StatusInternalServerError)
			return
		}
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
	if idUbicacionStr := r.FormValue("id_departamento"); idUbicacionStr != "" {
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

	tx := db.GDB.Begin()
	if err := tx.Save(&agencia).Error; err != nil {
		http.Error(w, "Error al actualizar la agencia", http.StatusInternalServerError)
		return
	}
	files := r.MultipartForm.File["fotos[]"]
	if len(files) > 0 {
		if err := tx.Where("id_agencia = ?", agencia.ID).Delete(&models.FotosAgencia{}).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al eliminar fotos anteriores", http.StatusInternalServerError)
			return
		}
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				tx.Rollback()
				http.Error(w, "Error al abrir una foto: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			ext := filepath.Ext(fileHeader.Filename)
			OrdenFoto, err := strconv.Atoi(strings.TrimSuffix((fileHeader.Filename), ext))
			if err != nil {
				tx.Rollback()
				http.Error(w, "Error al procesar el orden de la foto", http.StatusBadRequest)
				return
			}

			nombreFoto := fmt.Sprintf("foto_agencia_%s%s", uuid.New().String(), ext)
			rutaFoto := "internal/images/agencias/" + nombreFoto

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

			foto := models.FotosAgencia{
				IdAgencia: agencia.ID,
				Foto:      rutaFoto,
				Orden:     uint(OrdenFoto),
			}

			if err := tx.Create(&foto).Error; err != nil {
				tx.Rollback()
				http.Error(w, "Error al guardar la foto en base de datos", http.StatusInternalServerError)
				return
			}
		}
	}

	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agencia)
}

// Steps
func ObtenerAgenciaDatosGenerales(w http.ResponseWriter, r *http.Request) {
	id_agencia := mux.Vars(r)["id"]
	var datosGenerales types.AgenciaDatosGenerales

	query := `
		SELECT 
			nombre,
			telefono,
			correo,
			id_encargado,
			descripcion,
			id_departamento,
			direccion,
			estado
		FROM "GestAgencias"
		WHERE id_agencia = $1`

	if err := db.GDB.Raw(query, id_agencia).Scan(&datosGenerales).Error; err != nil {
		http.Error(w, "Error al obtener datos generales de la agencia", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(datosGenerales)
}

func ObtenerAgenciaFotos(w http.ResponseWriter, r *http.Request) {
	id_agencia := mux.Vars(r)["id"]
	var fotos []types.Foto

	query := `
		SELECT 
			id_foto as id,
			foto,
			orden
		FROM fotos_agencias
		WHERE id_agencia = $1
		ORDER BY orden ASC`

	if err := db.GDB.Raw(query, id_agencia).Scan(&fotos).Error; err != nil {
		http.Error(w, "Error al obtener fotos de la agencia", http.StatusInternalServerError)
		return
	}

	for i, foto := range fotos {
		base64img, err := encodeImageToBase64(foto.Foto)
		if err != nil {
			fmt.Printf("Error codificando una foto")

		}
		fotos[i].Foto = base64img
	}

	agenciaFotos := types.AgenciaFotos{
		Fotos: datatypes.NewJSONType(fotos),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agenciaFotos)
}

// Mod Steps
func ModificarAgenciaDatosGenerales(w http.ResponseWriter, r *http.Request) {
	id_agencia := mux.Vars(r)["id"]

	var agenciaExistente models.Agencia
	err := db.GDB.
		Where("id_agencia = ?", id_agencia).
		First(&agenciaExistente).
		Error
	if err != nil {
		http.Error(w, "Agencia no encontrada", http.StatusNotFound)
		return
	}

	var nuevosDatos types.AgenciaDatosGenerales
	if err := json.NewDecoder(r.Body).Decode(&nuevosDatos); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	agenciaExistente.Nombre = nuevosDatos.Nombre
	agenciaExistente.Telefono = nuevosDatos.Telefono
	agenciaExistente.Correo = nuevosDatos.Correo
	agenciaExistente.IdEncargado = nuevosDatos.IdEncargado
	agenciaExistente.Descripcion = nuevosDatos.Descripcion
	agenciaExistente.IdDepartamento = nuevosDatos.IdDepartamento
	agenciaExistente.Direccion = nuevosDatos.Direccion
	agenciaExistente.Estado = nuevosDatos.Estado

	if err := db.GDB.Save(&agenciaExistente).Error; err != nil {
		http.Error(w, "Error al modificar datos generales de la agencia", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agenciaExistente)
}

func ModificarAgenciaFotos(w http.ResponseWriter, r *http.Request) {
	id_agencia := mux.Vars(r)["id"]

	err := r.ParseMultipartForm(50 << 20) // 50 MB
	if err != nil {
		http.Error(w, "Error al parsear el formulario: "+err.Error(), http.StatusBadRequest)
		return
	}

	var agenciaExistente models.Agencia
	if err := db.GDB.Where("id_agencia = ?", id_agencia).First(&agenciaExistente).Error; err != nil {
		http.Error(w, "Agencia no encontrada", http.StatusNotFound)
		return
	}

	files := r.MultipartForm.File["fotos[]"]
	if len(files) == 0 {
		http.Error(w, "No se recibieron fotos", http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()

	var fotosAnteriores []models.FotosAgencia
	if err := tx.Where("id_agencia = ?", id_agencia).Find(&fotosAnteriores).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al obtener fotos anteriores", http.StatusInternalServerError)
		return
	}

	for _, foto := range fotosAnteriores {
		if foto.Foto != "" && foto.Foto != "N/A" {
			_ = os.Remove(foto.Foto) // Ignorar error si no existe
		}
	}

	if err := tx.Where("id_agencia = ?", id_agencia).Delete(&models.FotosAgencia{}).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al eliminar fotos anteriores", http.StatusInternalServerError)
		return
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al abrir una foto: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		ext := filepath.Ext(fileHeader.Filename)
		OrdenFoto, err := strconv.Atoi(strings.TrimSuffix(fileHeader.Filename, ext))
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al procesar el orden de la foto", http.StatusBadRequest)
			return
		}

		nombreFoto := fmt.Sprintf("foto_agencia_%s%s", uuid.New().String(), ext)
		rutaFoto := "internal/images/agencias/" + nombreFoto

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

		foto := models.FotosAgencia{
			IdAgencia: agenciaExistente.ID,
			Foto:      rutaFoto,
			Orden:     uint(OrdenFoto),
		}

		if err := tx.Create(&foto).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al guardar la foto en base de datos", http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Fotos actualizadas exitosamente"})
}
