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

func ObtenerAtraccionesFotos(w http.ResponseWriter, r *http.Request) {
	var atracciones []types.AtraccionesFotos

	err := db.GDB.Raw(services.QueryAtraccionesFoto).Scan(&atracciones).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	for i := range atracciones {
		fotoBase64 := ""
		if atracciones[i].Foto != "N/A" {
			if encoded, err := encodeImageToBase64(atracciones[i].Foto); err == nil {
				fotoBase64 = encoded
			}
		}
		if fotoBase64 == "" {
			fotoBase64 = "N/A"
		}
		atracciones[i].Foto = fotoBase64
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&atracciones)

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

	fmt.Printf("direccion %s", atraccion.Direccion)

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

	subcatsStr := r.FormValue("subcategorias")
	subcatsSlice := strings.Split(subcatsStr, ",")

	for _, s := range subcatsSlice {
		idSubcat, err := strconv.ParseUint(strings.TrimSpace(s), 10, 64)
		if err != nil {
			tx.Rollback()
			http.Error(w, "subcategoría inválida: "+s, http.StatusBadRequest)
			return
		}

		relacion := models.AtraccionesCategorias{
			IDAtraccion:    nuevaAtraccion.ID,
			IDSubCategoria: uint(idSubcat),
		}

		if err := tx.Create(&relacion).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al guardar la relación atracción-subcategoría", http.StatusInternalServerError)
			return
		}
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
		http.Error(w, "Atracción no encontrada", http.StatusNotFound)
		return
	}

	err := r.ParseMultipartForm(50 << 20) // 50 MB
	if err != nil {
		http.Error(w, "Error al parsear el formulario: "+err.Error(), http.StatusBadRequest)
		return
	}

	precio, _ := strconv.ParseFloat(r.FormValue("precio"), 64)
	idUbicacion, _ := strconv.ParseUint(r.FormValue("id_ubicacion"), 10, 64)
	idEncargado, _ := strconv.ParseUint(r.FormValue("id_encargado"), 10, 64)

	atraccion.Nombre = r.FormValue("nombre")
	atraccion.Direccion = r.FormValue("direccion")
	atraccion.Descripcion = r.FormValue("descripcion")
	atraccion.HorarioApertura = r.FormValue("horario_apertura")
	atraccion.HorarioCierre = r.FormValue("horario_cierre")
	atraccion.Precio = precio
	atraccion.Estado = r.FormValue("estado") == "true"
	atraccion.IdEncargado = uint(idEncargado)
	atraccion.IdUbicacion = uint(idUbicacion)

	tx := db.GDB.Begin()
	if err := tx.Save(&atraccion).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al actualizar la atracción", http.StatusInternalServerError)
		return
	}
	files := r.MultipartForm.File["fotos[]"]
	if len(files) > 0 {
		if err := tx.Where("id_atraccion = ?", atraccion.ID).Delete(&models.FotosAtracciones{}).Error; err != nil {
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

			nombreFoto := fmt.Sprintf("foto_atraccion_%s%s", uuid.New().String(), ext)
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
				IdAtraccion: atraccion.ID,
				Foto:        rutaFoto,
				Orden:       uint(OrdenFoto),
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

// Steps
func ObtenerAtraccionDatosGenerales(w http.ResponseWriter, r *http.Request) {
	id_atraccion := mux.Vars(r)["id"]

	var datosGenerales types.AtraccionDatosGenerales

	query := `
		SELECT 
			gat.nombre,
			gat.id_encargado,
			gat.direccion,
			gat.descripcion,
			gat.id_ubicacion,
			p.id_departamento
		FROM "GestAtraccionesTuristicas" gat
		LEFT JOIN provincias p ON gat.id_ubicacion = p.id_provincia
		WHERE id_atraccion = $1`

	if err := db.GDB.Raw(query, id_atraccion).Scan(&datosGenerales).Error; err != nil {
		http.Error(w, "Error al obtener datos generales de la atracción", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(datosGenerales)
}

func ObtenerAtraccionDatosEspecificos(w http.ResponseWriter, r *http.Request) {
	id_atraccion := mux.Vars(r)["id"]

	var datosEspecificos struct {
		HorarioApertura string  `json:"horario_apertura"`
		HorarioCierre   string  `json:"horario_cierre"`
		Precio          float64 `json:"precio"`
		Estado          bool    `json:"estado"`
	}

	query := `
		SELECT 
			horario_apertura,
			horario_cierre,
			precio,
			estado
		FROM "GestAtraccionesTuristicas"
		WHERE id_atraccion = $1`

	if err := db.GDB.Raw(query, id_atraccion).Scan(&datosEspecificos).Error; err != nil {
		http.Error(w, "Error al obtener datos específicos de la atracción", http.StatusInternalServerError)
		return
	}

	// Obtener subcategorías
	var subcategorias []types.SubCategoria

	querySubcategorias := `
		SELECT 
			s.id_subcategoria as id,
			s.nombre,
			s.descripcion
		FROM subcategorias s
		INNER JOIN atracciones_categorias ac ON s.id_subcategoria = ac.id_subcategoria
		WHERE ac.id_atraccion = $1 AND s.estado = true
		ORDER BY s.nombre ASC`

	if err := db.GDB.Raw(querySubcategorias, id_atraccion).Scan(&subcategorias).Error; err != nil {
		http.Error(w, "Error al obtener subcategorías de la atracción", http.StatusInternalServerError)
		return
	}

	resultado := types.AtraccionDatoEspecificos{
		HorarioApertura: datosEspecificos.HorarioApertura,
		HorarioCierre:   datosEspecificos.HorarioCierre,
		Precio:          datosEspecificos.Precio,
		Estado:          datosEspecificos.Estado,
		Categorias:      datatypes.NewJSONType(subcategorias),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}

func ObtenerAtraccionFotos(w http.ResponseWriter, r *http.Request) {
	id_atraccion := mux.Vars(r)["id"]

	var fotos []types.Foto

	query := `
		SELECT 
			id_foto as id,
			foto,
			orden
		FROM fotos_atracciones
		WHERE id_atraccion = $1
		ORDER BY orden ASC`

	if err := db.GDB.Raw(query, id_atraccion).Scan(&fotos).Error; err != nil {
		http.Error(w, "Error al obtener fotos de la atracción", http.StatusInternalServerError)
		return
	}

	for i, foto := range fotos {
		base64img, err := encodeImageToBase64(foto.Foto)
		if err != nil {
			fmt.Printf("Error codificando una foto")

		}
		fotos[i].Foto = base64img
	}

	atraccionFotos := types.AtraccionFotos{
		Fotos: datatypes.NewJSONType(fotos),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(atraccionFotos)
}

// Mod Steps
func ModificarAtraccionDatosGenerales(w http.ResponseWriter, r *http.Request) {
	id_atraccion := mux.Vars(r)["id"]

	var atraccionExistente models.AtraccionTuristica
	err := db.GDB.
		Where("id_atraccion = ?", id_atraccion).
		First(&atraccionExistente).
		Error
	if err != nil {
		http.Error(w, "Atracción no encontrada", http.StatusNotFound)
		return
	}

	var nuevosDatos types.AtraccionDatosGenerales
	if err := json.NewDecoder(r.Body).Decode(&nuevosDatos); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	atraccionExistente.Nombre = nuevosDatos.Nombre
	atraccionExistente.IdEncargado = nuevosDatos.IDEncargado
	atraccionExistente.Direccion = nuevosDatos.Direccion
	atraccionExistente.IdUbicacion = nuevosDatos.IDUbicacion
	atraccionExistente.Descripcion = nuevosDatos.Descripcion

	if err := db.GDB.Save(&atraccionExistente).Error; err != nil {
		http.Error(w, "Error al modificar datos generales de la atracción", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(atraccionExistente)
}

func ModificarAtraccionDatosEspecificos(w http.ResponseWriter, r *http.Request) {
	id_atraccion := mux.Vars(r)["id"]

	var atraccionExistente models.AtraccionTuristica
	err := db.GDB.
		Where("id_atraccion = ?", id_atraccion).
		First(&atraccionExistente).
		Error
	if err != nil {
		http.Error(w, "Atracción no encontrada", http.StatusNotFound)
		return
	}

	var nuevosDatos struct {
		HorarioApertura string                                   `json:"horario_apertura"`
		HorarioCierre   string                                   `json:"horario_cierre"`
		Precio          float64                                  `json:"precio"`
		Estado          bool                                     `json:"estado"`
		Categorias      datatypes.JSONType[[]types.SubCategoria] `json:"subcategorias"`
	}
	if err := json.NewDecoder(r.Body).Decode(&nuevosDatos); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	atraccionExistente.HorarioApertura = nuevosDatos.HorarioApertura
	atraccionExistente.HorarioCierre = nuevosDatos.HorarioCierre
	atraccionExistente.Precio = nuevosDatos.Precio
	atraccionExistente.Estado = nuevosDatos.Estado

	if err := db.GDB.Save(&atraccionExistente).Error; err != nil {
		http.Error(w, "Error al modificar datos específicos de la atracción", http.StatusInternalServerError)
		return
	}

	// Actualizar subcategorías
	// Primero eliminar las relaciones existentes
	if err := db.GDB.Exec("DELETE FROM atracciones_categorias WHERE id_atraccion = ?", id_atraccion).Error; err != nil {
		http.Error(w, "Error al actualizar subcategorías", http.StatusInternalServerError)
		return
	}

	// Insertar las nuevas relaciones
	subcategorias := nuevosDatos.Categorias.Data()
	for _, subcat := range subcategorias {
		nuevaRelacion := models.AtraccionesCategorias{
			IDAtraccion:    atraccionExistente.ID,
			IDSubCategoria: subcat.ID,
		}
		if err := db.GDB.Create(&nuevaRelacion).Error; err != nil {
			http.Error(w, "Error al crear relaciones de subcategorías", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(atraccionExistente)
}
