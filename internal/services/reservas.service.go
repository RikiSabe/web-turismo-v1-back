package services

var QueryReservasUsuarioTODO = `
	SELECT 
		r.id_reserva as id, 
		r.fecha, 
		r.descripcion, 
		r.numero_personas, 
		r.estado
	from "GestReservas" as r
	WhERE r.id_usuario = ?`
