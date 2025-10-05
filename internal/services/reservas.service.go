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
		r.descripcion,
		r.estado,
		r.id_usuario,
		u.nombre AS nombre_usuario,
		r.id_paquete,
		p.nombre AS nombre_paquete
	FROM "GestReservas" as r
	JOIN "GestUsuarios" u ON r.id_usuario = u.id_usuario
	JOIN "GestPaquetesTuristicos" p ON r.id_paquete = p.id_paquete_turistico`

const QueryReservasByUsuario = `
	SELECT 
		COALESCE(
			jsonb_agg(
				jsonb_build_object(
					'id', r.id_reserva,
					'descripcion', r.descripcion,
					'numero_personas', r.numero_personas,
					'estado', r.estado,
					'paquete', jsonb_build_object(
						'nombre', pt.nombre,
						'descripcion', pt.descripcion
					)
				)
				ORDER BY r.created_at DESC
			),
			'[]'::jsonb
		) AS reservas
	FROM "GestReservas" r
	INNER JOIN "GestPaquetesTuristicos" pt ON r.id_paquete = pt.id_paquete_turistico
	WHERE r.id_usuario = ?`
