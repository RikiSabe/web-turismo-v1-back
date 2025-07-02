package services

var QueryPaqueteTuristicoTODO = `
	SELECT 
    p.id_paquete_turistico AS id,
    p.categoria,
    p.nombre,
    p.fecha,
    p.descripcion,
    p.precio,
    p.duracion,
    p.salida,
    p.estado,
    p.id_agencia,
    json_agg(
			DISTINCT jsonb_build_object(
				'id', a.id_atracciones,
				'tipo', a.tipo,
				'nombre', a.nombre,
				'ubicacion', a.ubicacion,
				'descripcion', a.descripcion,
				'horarios', a.horarios,
				'precio', a.precio,
				'estado', a.estado
			)
    ) FILTER (WHERE a.id_atracciones IS NOT NULL) AS atracciones
	FROM "GestPaquetesTuristicos" as p
	LEFT JOIN "PaquetesAtracciones" as pa ON p.id_paquete_turistico = pa.id_paquete
	LEFT JOIN "GestAtraccionesTuristicas" as a ON pa.id_atraccion = a.id_atracciones
	GROUP BY 
		p.id_paquete_turistico,
		p.categoria,
		p.nombre,
		p.fecha,
		p.descripcion,
		p.precio,
		p.duracion,
		p.salida,
		p.estado,
		p.id_agencia
	ORDER BY p.id_paquete_turistico DESC`

var QueryPaqueteTuristicoTODOByID = `
	SELECT 
		p.id_paquete_turistico AS id,
		p.categoria,
		p.nombre,
		p.fecha,
		p.descripcion,
		p.precio,
		p.duracion,
		p.salida,
		p.estado,
		p.id_agencia,
		json_agg(
			DISTINCT jsonb_build_object(
				'id', a.id_atracciones,
				'tipo', a.tipo,
				'nombre', a.nombre,
				'ubicacion', a.ubicacion,
				'descripcion', a.descripcion,
				'horarios', a.horarios,
				'precio', a.precio,
				'estado', a.estado
			)
		) FILTER (WHERE a.id_atracciones IS NOT NULL) AS atracciones
	FROM "GestPaquetesTuristicos" AS p
	LEFT JOIN "PaquetesAtracciones" AS pa ON p.id_paquete_turistico = pa.id_paquete
	LEFT JOIN "GestAtraccionesTuristicas" AS a ON pa.id_atraccion = a.id_atracciones
	WHERE p.id_paquete_turistico = ?
	GROUP BY 
		p.id_paquete_turistico,
		p.categoria,
		p.nombre,
		p.fecha,
		p.descripcion,
		p.precio,
		p.duracion,
		p.salida,
		p.estado,
		p.id_agencia
	ORDER BY p.id_paquete_turistico DESC`
