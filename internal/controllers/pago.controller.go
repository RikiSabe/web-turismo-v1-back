package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ReservaPagoDTO struct {
	IDReserva      uint   `json:"id_reserva"`
	Descripcion    string `json:"descripcion"`
	NumeroPersonas int    `json:"numero_personas"`
	EstadoReserva  bool   `json:"estado_reserva"`
	FechaReserva   string `json:"fecha_reserva"`

	IDUsuario       uint   `json:"id_usuario"`
	NombreUsuario   string `json:"nombre_usuario"`
	ApellidoUsuario string `json:"apellido_usuario"`
	CorreoUsuario   string `json:"correo_usuario"`
	TelefonoUsuario string `json:"telefono_usuario"`

	IDPaquete          uint    `json:"id_paquete"`
	NombrePaquete      string  `json:"nombre_paquete"`
	DescripcionPaquete string  `json:"descripcion_paquete"`
	ActividadesPaquete string  `json:"actividades_paquete"`
	TipoPaquete        string  `json:"tipo_paquete"`
	PrecioPaquete      float64 `json:"precio_paquete"`
	DuracionPaquete    string  `json:"duracion_paquete"`
	HoraInicial        string  `json:"hora_inicial"`
	HoraFinal          string  `json:"hora_final"`
	CupoTotal          uint    `json:"cupo_total"`
	CupoMinimo         uint    `json:"cupo_minimo"`

	IDAgencia        uint   `json:"id_agencia"`
	NombreAgencia    string `json:"nombre_agencia"`
	DireccionAgencia string `json:"direccion_agencia"`
	TelefonoAgencia  string `json:"telefono_agencia"`
	CorreoAgencia    string `json:"correo_agencia"`

	IDEncargado       uint   `json:"id_encargado"`
	NombreEncargado   string `json:"nombre_encargado"`
	ApellidoEncargado string `json:"apellido_encargado"`
	CorreoEncargado   string `json:"correo_encargado"`
	TelefonoEncargado string `json:"telefono_encargado"`
	FotoEncargado     string `json:"foto_encargado"`
}

var QueryReservasPago = `
	SELECT
		-- Datos de la reserva
		r.id_reserva,
		r.descripcion,
		r.numero_personas,
		r.estado AS estado_reserva,
		r.created_at AS fecha_reserva,
		
		-- Datos del usuario que reserva
		r.id_usuario,
		u.nombre AS nombre_usuario,
		CONCAT(u.apellido_paterno, ' ', COALESCE(u.apellido_materno, '')) AS apellido_usuario,
		u.correo AS correo_usuario,
		u.telefono AS telefono_usuario,
		
		-- Datos del paquete turístico
		p.id_paquete_turistico AS id_paquete,
		p.nombre AS nombre_paquete,
		p.descripcion AS descripcion_paquete,
		p.actividades AS actividades_paquete,
		p.tipo AS tipo_paquete,
		p.precio AS precio_paquete,
		p.duracion AS duracion_paquete,
		p.hora_inicial,
		p.hora_final,
		p.cupo_total,
		p.cupo_minimo,
		
		-- Datos de la agencia
		a.id_agencia,
		a.nombre AS nombre_agencia,
		a.direccion AS direccion_agencia,
		a.telefono AS telefono_agencia,
		a.correo AS correo_agencia,
		
		-- Datos del encargado de la agencia
		enc.id_usuario AS id_encargado,
		enc.nombre AS nombre_encargado,
		CONCAT(enc.apellido_paterno, ' ', COALESCE(enc.apellido_materno, '')) AS apellido_encargado,
		enc.correo AS correo_encargado,
		enc.telefono AS telefono_encargado,
		enc.foto AS foto_encargado

	FROM "GestReservas" r
	INNER JOIN "GestUsuarios" u ON r.id_usuario = u.id_usuario
	INNER JOIN "GestPaquetesTuristicos" p ON r.id_paquete = p.id_paquete_turistico
	INNER JOIN "GestAgencias" a ON p.id_agencia = a.id_agencia
	INNER JOIN "GestUsuarios" enc ON a.id_encargado = enc.id_usuario
	ORDER BY r.created_at DESC`

func ObtenerReservasPago(w http.ResponseWriter, r *http.Request) {
	var reservas []ReservaPagoDTO

	if err := db.GDB.Raw(QueryReservasPago).Scan(&reservas).Error; err != nil {
		http.Error(w, "Error al obtener reservas para pago", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservas)
}

type ReservaPagoUsuarioDTO struct {
	IDReserva      uint   `json:"id_reserva"`
	Descripcion    string `json:"descripcion"`
	NumeroPersonas int    `json:"numero_personas"`
	EstadoReserva  bool   `json:"estado_reserva"`
	FechaReserva   string `json:"fecha_reserva"`

	IDUsuario       uint   `json:"id_usuario"`
	NombreUsuario   string `json:"nombre_usuario"`
	ApellidoUsuario string `json:"apellido_usuario"`
	CorreoUsuario   string `json:"correo_usuario"`
	TelefonoUsuario string `json:"telefono_usuario"`

	IDPaquete          uint    `json:"id_paquete"`
	NombrePaquete      string  `json:"nombre_paquete"`
	DescripcionPaquete string  `json:"descripcion_paquete"`
	ActividadesPaquete string  `json:"actividades_paquete"`
	TipoPaquete        string  `json:"tipo_paquete"`
	PrecioPaquete      float64 `json:"precio_paquete"`
	DuracionPaquete    string  `json:"duracion_paquete"`
	HoraInicial        string  `json:"hora_inicial"`
	HoraFinal          string  `json:"hora_final"`
	CupoTotal          uint    `json:"cupo_total"`
	CupoMinimo         uint    `json:"cupo_minimo"`

	IDAgencia        uint   `json:"id_agencia"`
	NombreAgencia    string `json:"nombre_agencia"`
	DireccionAgencia string `json:"direccion_agencia"`
	TelefonoAgencia  string `json:"telefono_agencia"`
	CorreoAgencia    string `json:"correo_agencia"`

	IDEncargado       uint   `json:"id_encargado"`
	NombreEncargado   string `json:"nombre_encargado"`
	ApellidoEncargado string `json:"apellido_encargado"`
	CorreoEncargado   string `json:"correo_encargado"`
	TelefonoEncargado string `json:"telefono_encargado"`
	FotoEncargado     string `json:"foto_encargado"`
}

var QueryReservasPagoUsuario = `
	SELECT
		-- Datos de la reserva
		r.id_reserva,
		r.descripcion,
		r.numero_personas,
		r.estado AS estado_reserva,
		r.created_at AS fecha_reserva,
		
		-- Datos del usuario que reserva
		r.id_usuario,
		u.nombre AS nombre_usuario,
		CONCAT(u.apellido_paterno, ' ', COALESCE(u.apellido_materno, '')) AS apellido_usuario,
		u.correo AS correo_usuario,
		u.telefono AS telefono_usuario,
		
		-- Datos del paquete turístico
		p.id_paquete_turistico AS id_paquete,
		p.nombre AS nombre_paquete,
		p.descripcion AS descripcion_paquete,
		p.actividades AS actividades_paquete,
		p.tipo AS tipo_paquete,
		p.precio AS precio_paquete,
		p.duracion AS duracion_paquete,
		p.hora_inicial,
		p.hora_final,
		p.cupo_total,
		p.cupo_minimo,
		
		-- Datos de la agencia
		a.id_agencia,
		a.nombre AS nombre_agencia,
		a.direccion AS direccion_agencia,
		a.telefono AS telefono_agencia,
		a.correo AS correo_agencia,
		
		-- Datos del encargado de la agencia
		enc.id_usuario AS id_encargado,
		enc.nombre AS nombre_encargado,
		CONCAT(enc.apellido_paterno, ' ', COALESCE(enc.apellido_materno, '')) AS apellido_encargado,
		enc.correo AS correo_encargado,
		enc.telefono AS telefono_encargado,
		enc.foto AS foto_encargado

	FROM "GestReservas" r
	INNER JOIN "GestUsuarios" u ON r.id_usuario = u.id_usuario
	INNER JOIN "GestPaquetesTuristicos" p ON r.id_paquete = p.id_paquete_turistico
	INNER JOIN "GestAgencias" a ON p.id_agencia = a.id_agencia
	INNER JOIN "GestUsuarios" enc ON a.id_encargado = enc.id_usuario
	WHERE r.id_usuario = $1
	ORDER BY r.created_at DESC`

func ReservasPagoUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idUsuarioStr := vars["id"]

	idUsuario, err := strconv.ParseUint(idUsuarioStr, 10, 32)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	var reservas []ReservaPagoUsuarioDTO

	if err := db.GDB.Raw(QueryReservasPagoUsuario, idUsuario).Scan(&reservas).Error; err != nil {
		http.Error(w, "Error al obtener reservas del usuario", http.StatusInternalServerError)
		return
	}

	if reservas == nil {
		reservas = []ReservaPagoUsuarioDTO{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservas)
}

func EnviarComprobantePago(w http.ResponseWriter, r *http.Request) {
	id_reserva := mux.Vars(r)["id_reserva"]
	id_usuario := mux.Vars(r)["id_usuario"]

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error al parsear el formulario", http.StatusInternalServerError)
		return
	}

	direccionCoprobante := "N/A"

	file, handler, err := r.FormFile("comprobante")
	if err == nil {
		defer file.Close()

		nombreComprobante := fmt.Sprintf("foto_comprobante_%s%s", uuid.New().String(), filepath.Ext(handler.Filename))

		rutaFoto := "internal/images/comprobantes/" + nombreComprobante
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

		direccionCoprobante = rutaFoto
	} else if err != http.ErrMissingFile {
		http.Error(w, "Error al obtener la foto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	idReserva, err := strconv.ParseUint(id_reserva, 10, 64)
	if err != nil {
		http.Error(w, "id_reserva inválido", http.StatusBadRequest)
		return
	}

	idUsuario, err := strconv.ParseUint(id_usuario, 10, 64)
	if err != nil {
		http.Error(w, "id_usuario inválido", http.StatusBadRequest)
		return
	}

	newComprobante := models.Comprobante{
		IDReserva: uint(idReserva),
		IDUsuario: uint(idUsuario),
		Foto:      direccionCoprobante,
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&newComprobante).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al guardar el comprobante", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Comprobante enviado correctamente"})
}
