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

var QueryReservasTODO = `
	SELECT 
		r.id_reserva AS id,
		r.fecha,
		r.descripcion,
		r.numero_personas,
		r.estado,
		r.id_usuario,
		u.nombre AS nombre_usuario,
		r.id_paquete,
		p.nombre AS nombre_paquete
	FROM "GestReservas" as r
	JOIN "GestUsuarios" u ON r.id_usuario = u.id_usuario
	JOIN "GestPaquetesTuristicos" p ON r.id_paquete = p.id_paquete_turistico`
