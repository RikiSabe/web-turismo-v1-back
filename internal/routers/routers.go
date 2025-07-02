package routers

import (
	"net/http"
	c "web-turismo-v1/internal/controllers"

	"github.com/gorilla/mux"
)

func InitEndPoints(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()
	endPointsAPI(api)
}

func endPointsAPI(api *mux.Router) {
	// API
	v1 := api.PathPrefix("/v1").Subrouter()

	// Prefix
	v1AtraccionesTuristicas := v1.PathPrefix("/atracciones-turisticas").Subrouter()
	v1Agencias := v1.PathPrefix("/agencias").Subrouter()
	v1Roles := v1.PathPrefix("/roles").Subrouter()
	v1Usuarios := v1.PathPrefix("/usuarios").Subrouter()
	v1Reservas := v1.PathPrefix("/reservas").Subrouter()
	v1PaquetesTuristicos := v1.PathPrefix("/paquetes-turisticos").Subrouter()

	// Auth
	// v1
	v1.HandleFunc("/loginweb", c.Auth.AuthLoginWeb).Methods(http.MethodPost)

	// Atracciones turisticas
	// v1
	v1AtraccionesTuristicas.HandleFunc("/{id}", c.ObtenerAtraccionTuristica).Methods(http.MethodGet)
	v1AtraccionesTuristicas.HandleFunc("/{id}", c.ModificarAtraccionTuristica).Methods(http.MethodPut)
	v1AtraccionesTuristicas.HandleFunc("", c.ObtenerAtraccionesTuristicas).Methods(http.MethodGet)
	v1AtraccionesTuristicas.HandleFunc("", c.AgregarAtraccionTuristica).Methods(http.MethodPost)

	// Agencias
	// v1
	v1Agencias.HandleFunc("/{id}", c.ObtenerAgencia).Methods(http.MethodGet)
	v1Agencias.HandleFunc("/{id}", c.ModificarAgencia).Methods(http.MethodPut)
	v1Agencias.HandleFunc("", c.ObtenerAgencias).Methods(http.MethodGet)
	v1Agencias.HandleFunc("", c.AgregarAgencia).Methods(http.MethodPost)

	// Roles
	// v1
	v1Roles.HandleFunc("/{id}", c.ObtenerRol).Methods(http.MethodGet)
	v1Roles.HandleFunc("/{id}", c.ModificarRol).Methods(http.MethodPut)
	v1Roles.HandleFunc("", c.ObtenerRoles).Methods(http.MethodGet)
	v1Roles.HandleFunc("", c.AgregarRol).Methods(http.MethodPost)

	// Usuarios
	// v1
	v1Usuarios.HandleFunc("/{id}", c.ObtenerUsuario).Methods(http.MethodGet)
	v1Usuarios.HandleFunc("/{id}", c.ModificarUsuario).Methods(http.MethodPut)
	v1Usuarios.HandleFunc("", c.ObtenerUsuarios).Methods(http.MethodGet)
	v1Usuarios.HandleFunc("", c.AgregarUsuario).Methods(http.MethodPost)

	// Reservas
	// v1
	v1Reservas.HandleFunc("/usuario/{id}", c.ObtenerReservasUsuario).Methods(http.MethodGet)
	v1Reservas.HandleFunc("", c.HacerReserva).Methods(http.MethodPost)

	// Paquetes turisticos
	// v1
	v1PaquetesTuristicos.HandleFunc("", c.ObtenerPaquetesTuristicos).Methods(http.MethodGet)
	v1PaquetesTuristicos.HandleFunc("/{id}", c.ObtenerPaqueteTuristico).Methods(http.MethodGet)
	v1PaquetesTuristicos.HandleFunc("", c.CrearPaqueteTuristico).Methods(http.MethodPost)
}
