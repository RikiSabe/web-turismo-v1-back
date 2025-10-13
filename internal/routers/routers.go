package routers

import (
	"net/http"
	c "web-turismo-v1/internal/controllers"
	r "web-turismo-v1/internal/reports"

	"github.com/gorilla/mux"
)

func InitEndPoints(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()
	endPointsAPI(api)
	reportAPI(api)
}

func reportAPI(api *mux.Router) {
	v1 := api.PathPrefix("/v1").Subrouter()

	// Prefix
	// v1
	v1Reports := v1.PathPrefix("/reportes").Subrouter()
	v1Reports.HandleFunc("/usuarios", r.ReporteUsuarios).Methods(http.MethodGet)
	v1Reports.HandleFunc("/atracciones", r.ReporteAtracciones).Methods(http.MethodGet)
	v1Reports.HandleFunc("/agencias", r.ReporteAgencias).Methods(http.MethodGet)
}

func endPointsAPI(api *mux.Router) {
	// API
	v1 := api.PathPrefix("/v1").Subrouter()
	v2 := api.PathPrefix("/v2").Subrouter()

	// Prefix
	// v1
	v1Agencias := v1.PathPrefix("/agencias").Subrouter()
	v1Usuarios := v1.PathPrefix("/usuarios").Subrouter()
	v1Reservas := v1.PathPrefix("/reservas").Subrouter()
	v1PaquetesTuristicos := v1.PathPrefix("/paquetes-turisticos").Subrouter()
	v1Departamentos := v1.PathPrefix("/departamentos").Subrouter()
	v1Provincias := v1.PathPrefix("/provincias").Subrouter()
	v1Categorias := v1.PathPrefix("/categorias").Subrouter()
	v1SubCategorias := v1.PathPrefix("/subcategorias").Subrouter()
	v1Dias := v1.PathPrefix("/dias").Subrouter()
	v1Encargado := v1.PathPrefix("/encargado-qr").Subrouter()
	v1Pagos := v1.PathPrefix("/pagos").Subrouter()

	// v2
	v2Usuarios := v2.PathPrefix("/usuarios").Subrouter()
	v2AtraccionesTuristicas := v2.PathPrefix("/atracciones-turisticas").Subrouter()

	// Auth
	v1.HandleFunc("/loginweb", c.Auth.AuthLoginWeb).Methods(http.MethodPost)
	v1.HandleFunc("/registro", c.Auth.AuthRegisterWeb).Methods(http.MethodPost)

	// Atracciones turisticas
	v2AtraccionesTuristicas.HandleFunc("/encargado/{id}", c.ObtenerEncargadoAtraccionTuristica).Methods(http.MethodGet)
	v2AtraccionesTuristicas.HandleFunc("/fotos", c.ObtenerAtraccionesFotos).Methods(http.MethodGet)
	v2AtraccionesTuristicas.HandleFunc("/datos-generales/{id}", c.ObtenerAtraccionDatosGenerales).Methods(http.MethodGet)
	v2AtraccionesTuristicas.HandleFunc("/datos-generales/{id}", c.ModificarAtraccionDatosGenerales).Methods(http.MethodPut)
	v2AtraccionesTuristicas.HandleFunc("/datos-especificos/{id}", c.ObtenerAtraccionDatosEspecificos).Methods(http.MethodGet)
	v2AtraccionesTuristicas.HandleFunc("/datos-especificos/{id}", c.ModificarAtraccionDatosEspecificos).Methods(http.MethodPut)
	v2AtraccionesTuristicas.HandleFunc("/fotos/{id}", c.ObtenerAtraccionFotos).Methods(http.MethodGet)
	v2AtraccionesTuristicas.HandleFunc("/{id}", c.ObtenerAtraccionTuristica).Methods(http.MethodGet)
	v2AtraccionesTuristicas.HandleFunc("/{id}", c.ModificarAtraccionTuristica).Methods(http.MethodPut)
	v2AtraccionesTuristicas.HandleFunc("", c.ObtenerAtraccionesTuristicas).Methods(http.MethodGet)
	v2AtraccionesTuristicas.HandleFunc("", c.AgregarAtraccionTuristica).Methods(http.MethodPost)

	// Agencias
	v1Agencias.HandleFunc("/{id}", c.ObtenerAgencia).Methods(http.MethodGet)
	v1Agencias.HandleFunc("/{id}", c.ModificarAgencia).Methods(http.MethodPut)
	v1Agencias.HandleFunc("/datos-generales/{id}", c.ObtenerAgenciaDatosGenerales).Methods(http.MethodGet)
	v1Agencias.HandleFunc("/datos-generales/{id}", c.ModificarAgenciaDatosGenerales).Methods(http.MethodPut)
	v1Agencias.HandleFunc("/fotos/{id}", c.ObtenerAgenciaFotos).Methods(http.MethodGet)
	v1Agencias.HandleFunc("", c.ObtenerAgencias).Methods(http.MethodGet)
	v1Agencias.HandleFunc("", c.AgregarAgencia).Methods(http.MethodPost)

	// Usuarios
	v2Usuarios.HandleFunc("/menu/{id}", c.ObtenerUsuarioMenu).Methods(http.MethodGet)
	v1Usuarios.HandleFunc("/datos-personales/{id}", c.ObtenerUsuarioDatosPersonales).Methods(http.MethodGet)
	v1Usuarios.HandleFunc("/datos-personales/{id}", c.ModificarUsuarioDatosPersonales).Methods(http.MethodPut)
	v1Usuarios.HandleFunc("/datos-privados/{id}", c.ObtenerUsuarioDatosPrivados).Methods(http.MethodGet)
	v1Usuarios.HandleFunc("/datos-privados/{id}", c.ModificarUsuarioDatosPrivados).Methods(http.MethodPut)
	v1Usuarios.HandleFunc("/foto/{id}", c.ObtenerUsuarioFoto).Methods(http.MethodGet)
	v2Usuarios.HandleFunc("/{id}", c.ObtenerUsuario).Methods(http.MethodGet)
	v2Usuarios.HandleFunc("/{id}", c.ModificarUsuario).Methods(http.MethodPut)
	v1Usuarios.HandleFunc("", c.ObtenerUsuarios).Methods(http.MethodGet)
	v2Usuarios.HandleFunc("", c.AgregarUsuario).Methods(http.MethodPost)

	// Reservas
	v1Reservas.HandleFunc("/usuario/{id}", c.ObtenerReservasPorUsuario).Methods(http.MethodGet)
	v1Reservas.HandleFunc("/decision/{id}", c.DecisionReserva).Methods(http.MethodPost)
	v1Reservas.HandleFunc("/{id}", c.ObtenerReserva).Methods(http.MethodGet)
	v1Reservas.HandleFunc("", c.HacerReserva).Methods(http.MethodPost)
	v1Reservas.HandleFunc("", c.ObtenerReservas).Methods(http.MethodGet)

	// Paquetes turisticos
	v1PaquetesTuristicos.HandleFunc("/foto", c.ObtenerPaquetesTuristicosFoto).Methods(http.MethodGet)
	v1PaquetesTuristicos.HandleFunc("/foto/{id}", c.ObtenerPaqueteTuristicoFoto).Methods(http.MethodGet)
	v1PaquetesTuristicos.HandleFunc("", c.ObtenerPaquetesTuristicos).Methods(http.MethodGet)
	v1PaquetesTuristicos.HandleFunc("", c.CrearPaqueteTuristico).Methods(http.MethodPost)

	// Departamentos
	v1Departamentos.HandleFunc("", c.ObtenerDepartamentos).Methods(http.MethodGet)

	// Provincias
	v1Provincias.HandleFunc("/{id}", c.ObtenerProvinciasbyDepartamento).Methods(http.MethodGet)

	// Categorias
	v1Categorias.HandleFunc("/{id}", c.ObtenerCategoria).Methods(http.MethodGet)
	v1Categorias.HandleFunc("/{id}", c.ModificarCategoria).Methods(http.MethodPut)
	v1Categorias.HandleFunc("", c.ObtenerCategorias).Methods(http.MethodGet)
	v1Categorias.HandleFunc("", c.AgregarCategoria).Methods(http.MethodPost)

	// Sub Categorias
	v1SubCategorias.HandleFunc("/{id}", c.ObtenerSubCategoriasPorCategoria).Methods(http.MethodGet)

	// Dias
	v1Dias.HandleFunc("", c.ObtenerDias).Methods(http.MethodGet)

	// Encargado QR
	v1Encargado.HandleFunc("/{id}", c.AsignarQR).Methods(http.MethodPost)
	v1Encargado.HandleFunc("/{id}", c.ObtenerQR).Methods(http.MethodGet)

	// Pagos
	v1Pagos.HandleFunc("/comprobante/{id_usuario}/{id_reserva}", c.EnviarComprobantePago).Methods(http.MethodPost)
	v1Pagos.HandleFunc("/{id}", c.ReservasPagoUsuario).Methods(http.MethodGet)
	v1Pagos.HandleFunc("", c.ObtenerReservasPago).Methods(http.MethodGet)
	v1Pagos.HandleFunc("/comprobante/{id_usuario}/{id_reserva}", c.ObtenerDetallesComprobante).Methods(http.MethodGet)
}
