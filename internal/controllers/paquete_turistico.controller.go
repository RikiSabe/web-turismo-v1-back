package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/models"
	"web-turismo-v1/internal/services"
	"web-turismo-v1/internal/types"

	"github.com/gorilla/mux"
)

func CrearPaqueteTuristico(w http.ResponseWriter, r *http.Request) {

	var bodyStruct struct {
		Nombre      string `json:"nombre"`
		IDAgencia   uint   `json:"id_agencia"`
		Descripcion string `json:"descripcion"`
		Actividades string `json:"actividades"`
		Visible     string `json:"visible"`
		Precio      string `json:"precio"`

		Tipo             string   `json:"tipo"`
		DiasConcurrentes []uint   `json:"diasconcurrentes"`
		FechaRango       []string `json:"fecharango"`
		FechaUnica       string   `json:"fechaunica"`
		HoraInicial      string   `json:"hora_inicial"`
		CupoTotal        string   `json:"cupo_total"`
		CupoMinimo       string   `json:"cupo_minimo"`

		Atracciones string `json:"atracciones"` //Volver a variables (viene en string)
	}

	if err := json.NewDecoder(r.Body).Decode(&bodyStruct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, err := json.MarshalIndent(bodyStruct, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("body (JSON):")
	fmt.Println(string(jsonData))

	precio, err := strconv.ParseFloat(bodyStruct.Precio, 64)
	if err != nil {
		http.Error(w, "Precio inválido", http.StatusBadRequest)
		return
	}

	cupoTotal, err := strconv.ParseUint(bodyStruct.CupoTotal, 10, 64)
	if err != nil {
		http.Error(w, "Cupo total inválido", http.StatusBadRequest)
		return
	}

	cupoMinimo, err := strconv.ParseUint(bodyStruct.CupoMinimo, 10, 64)
	if err != nil {
		http.Error(w, "Cupo mínimo inválido", http.StatusBadRequest)
		return
	}

	// convertir visible en bool
	visible := false
	if strings.ToLower(bodyStruct.Visible) == "visible" || bodyStruct.Visible == "true" {
		visible = true
	}

	paqueteTuristico := models.PaqueteTuristico{
		Nombre:      bodyStruct.Nombre,
		IDAgencia:   bodyStruct.IDAgencia,
		Descripcion: bodyStruct.Descripcion,
		Actividades: bodyStruct.Actividades,
		Tipo:        bodyStruct.Tipo,
		Precio:      precio,
		Duracion:    "0",
		HoraInicial: bodyStruct.HoraInicial,
		HoraFinal:   "N/A", // lo puedes definir más adelante si lo calculas
		CupoTotal:   uint(cupoTotal),
		CupoMinimo:  uint(cupoMinimo),
		Visible:     visible,
		Estado:      true,
		Promocional: false,
	}

	type AtraccionData struct {
		ID       uint   `json:"id"`
		Nombre   string `json:"nombre"`
		Orden    uint   `json:"orden"`
		Duracion uint   `json:"duracion"`
	}

	// Falta poner los tipos
	tx := db.GDB.Begin()
	if err := tx.Create(&paqueteTuristico).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Paquete Turistico", http.StatusInternalServerError)
		return
	}

	var atraccionesData []AtraccionData
	if err := json.Unmarshal([]byte(bodyStruct.Atracciones), &atraccionesData); err != nil {
		http.Error(w, "Error al parsear atracciones: "+err.Error(), http.StatusBadRequest)
		return
	}

	if bodyStruct.Tipo == "Días concurrentes" {

		for _, idxDia := range bodyStruct.DiasConcurrentes {
			data := models.VigenciaDiasConcurrentes{
				IDPaquete: paqueteTuristico.ID,
				IDDia:     idxDia,
			}
			if err := tx.Create(&data).Error; err != nil {
				http.Error(w, "Error al agregar tipo DC", http.StatusInternalServerError)
				return
			}
		}
	}

	if bodyStruct.Tipo == "Rango de días especificos" {
		fechaInicioParsed, err := time.Parse("2006-01-02", bodyStruct.FechaRango[0])
		if err != nil {
			tx.Rollback()
			http.Error(w, "Formato de fecha inválido para FechaInicio, se espera YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		fechaFinParsed, err := time.Parse("2006-01-02", bodyStruct.FechaRango[1])
		if err != nil {
			tx.Rollback()
			http.Error(w, "Formato de fecha inválido para FechaFin, se espera YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		data := models.VigenciaRango{
			IDPaquete:   paqueteTuristico.ID,
			FechaInicio: fechaInicioParsed,
			FechaFin:    fechaFinParsed,
		}
		if err := tx.Create(&data).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al agregar tipo RD", http.StatusInternalServerError)
			return
		}
	}

	if bodyStruct.Tipo == "Único día" {
		fechaUnicaParsed, err := time.Parse("2006-01-02", bodyStruct.FechaUnica)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Formato de fecha inválido para fechaunica, se espera YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		data := models.VigenciaUnica{
			IDPaquete: paqueteTuristico.ID,
			Fecha:     fechaUnicaParsed,
		}
		if err := tx.Create(&data).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al agregar tipo UD", http.StatusInternalServerError)
			return
		}
	}

	// --- Debug ---
	fmt.Println("Atracciones parseadas:")
	for _, a := range atraccionesData {
		fmt.Printf("  - ID: %d, Nombre: %s, Orden: %d, Duración: %d\n", a.ID, a.Nombre, a.Orden, a.Duracion)
	}

	for _, atr := range atraccionesData {
		paqueteAtr := models.PaqueteAtraccion{
			IDPaquete:   paqueteTuristico.ID,
			IDAtraccion: atr.ID,
			Orden:       atr.Orden,
			Duracion:    atr.Duracion,
		}

		if err := tx.Create(&paqueteAtr).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al agregar atracción al paquete: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&bodyStruct)
}

func ObtenerPaquetesTuristicos(w http.ResponseWriter, r *http.Request) {
	var paquetes []types.PaqueteTuristicoTODO

	if err := db.GDB.Raw(services.QueryPaqueteTuristicoTODO).Scan(&paquetes).Error; err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(paquetes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ObtenerPaqueteTuristico(w http.ResponseWriter, r *http.Request) {
	id_paquete := mux.Vars(r)["id"]

	var paquete types.PaqueteTuristicoTODO
	if err := db.GDB.Raw(services.QueryPaqueteTuristicoTODOByID, id_paquete).Scan(&paquete).Error; err != nil {
		http.Error(w, "Error al obtener el paquete turístico", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(paquete); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ObtenerPaquetesTuristicosFoto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var paquetes []models.PaqueteTuristico
	if err := db.GDB.Find(&paquetes).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type FotoAtraccion struct {
		IDAtraccion uint   `json:"id_atraccion"`
		Nombre      string `json:"nombre"`
		Foto        string `json:"foto"`
	}

	type PaqueteResponse struct {
		models.PaqueteTuristico
		Atracciones []FotoAtraccion `json:"atracciones"`
		Vigencias   any             `json:"vigencias,omitempty"`
	}

	var paquetesResp []PaqueteResponse

	for _, paquete := range paquetes {
		// Obtener atracciones del paquete
		var relaciones []models.PaqueteAtraccion
		if err := db.GDB.
			Where("id_paquete = ?", paquete.ID).
			Preload("Atraccion").
			Find(&relaciones).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var atraccionesResp []FotoAtraccion

		for _, rel := range relaciones {
			// Obtener la primera foto (orden = 1)
			var foto models.FotosAtracciones
			db.GDB.
				Where("id_atraccion = ? AND orden = 1", rel.IDAtraccion).
				First(&foto)

			encoded := ""
			if encodedfunc, err := encodeImageToBase64(foto.Foto); err == nil {
				encoded = encodedfunc
			}

			atraccionesResp = append(atraccionesResp, FotoAtraccion{
				IDAtraccion: rel.Atraccion.ID,
				Nombre:      rel.Atraccion.Nombre,
				Foto:        encoded,
			})
		}

		// encodear la foto encodeImageToBase64(foto)

		// Obtener vigencias según tipo
		var vigencias any
		switch paquete.Tipo {
		case "Días concurrentes":
			var dias []models.VigenciaDiasConcurrentes
			db.GDB.Where("id_paquete = ?", paquete.ID).Find(&dias)
			vigencias = dias
		case "Rango de días especificos":
			var rangos []models.VigenciaRango
			db.GDB.Where("id_paquete = ?", paquete.ID).Find(&rangos)
			vigencias = rangos
		case "Único día":
			var unica []models.VigenciaUnica
			db.GDB.Where("id_paquete = ?", paquete.ID).Find(&unica)
			vigencias = unica
		}

		paquetesResp = append(paquetesResp, PaqueteResponse{
			PaqueteTuristico: paquete,
			Atracciones:      atraccionesResp,
			Vigencias:        vigencias,
		})
	}

	if err := json.NewEncoder(w).Encode(paquetesResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ObtenerPaqueteTuristicoFoto(w http.ResponseWriter, r *http.Request) {
	id_paquete := mux.Vars(r)["id"]

	type AtraccionDetalle struct {
		ID          uint   `json:"id"`
		Nombre      string `json:"nombre"`
		Foto        string `json:"foto"`
		Duracion    uint   `json:"duracion"`
		Descripcion string `json:"descripcion"`
		Orden       uint   `json:"orden"`
	}

	type PaqueteTuristicoTODO struct {
		ID            uint               `json:"id"`
		Nombre        string             `json:"nombre"`
		Descripcion   string             `json:"descripcion"`
		Precio        float64            `json:"precio"`
		Subcategorias []string           `json:"subcategorias"`
		HoraInicio    string             `json:"hora_inicio"`
		Atracciones   []AtraccionDetalle `json:"atracciones"`
	}

	var result struct {
		Paquete json.RawMessage `gorm:"column:paquete"`
	}

	if err := db.GDB.Raw(services.QueryPaqueteTuristicoUnique, id_paquete).Scan(&result).Error; err != nil {
		http.Error(w, "Error al obtener el paquete turístico", http.StatusInternalServerError)
		return
	}

	var paquete PaqueteTuristicoTODO
	if err := json.Unmarshal(result.Paquete, &paquete); err != nil {
		http.Error(w, "Error al procesar el paquete turístico", http.StatusInternalServerError)
		return
	}

	// Encodear las fotos de cada atracción
	for i := range paquete.Atracciones {
		encoded := ""
		if paquete.Atracciones[i].Foto != "" && paquete.Atracciones[i].Foto != "N/A" {
			if encodedfunc, err := encodeImageToBase64(paquete.Atracciones[i].Foto); err == nil {
				encoded = encodedfunc
			}
		}
		paquete.Atracciones[i].Foto = encoded
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(paquete); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
