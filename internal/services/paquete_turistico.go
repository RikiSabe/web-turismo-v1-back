package services

var QueryPaqueteTuristicoTODO = `
	SELECT 
		gpt.id_paquete_turistico AS id, 
		gpt.nombre, 
		gpt.descripcion, 
		gpt.precio, 
		gpt.hora_inicial, 
		gpt.estado, 
		gpt.visible
	FROM "GestPaquetesTuristicos" AS gpt
	ORDER BY gpt.id_paquete_turistico;`

var QueryPaqueteTuristicoFoto = `
	`

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
				'id', a.id_atraccion,
				'tipo', a.categoria,
				'nombre', a.nombre,
				'ubicacion', a.direccion,
				'descripcion', a.descripcion,
				'horario_apertura', a.horario_apertura,
				'horario_cierre', a.horario_cierre,
				'precio', a.precio,
				'estado', a.estado
			)
		) FILTER (WHERE a.id_atraccion IS NOT NULL) AS atracciones
	FROM "GestPaquetesTuristicos" AS p
	LEFT JOIN "PaquetesAtracciones" AS pa ON p.id_paquete_turistico = pa.id_paquete
	LEFT JOIN "GestAtraccionesTuristicas" AS a ON pa.id_atraccion = a.id_atraccion
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

var QueryPaqueteTuristicoUnique = `
	SELECT 
		jsonb_build_object(
			'id', pt.id_paquete_turistico,
			'nombre', pt.nombre,
			'descripcion', pt.descripcion,
			'precio', pt.precio,
			'subcategorias', COALESCE(
				(
					SELECT jsonb_agg(DISTINCT sc.nombre ORDER BY sc.nombre)
					FROM "PaquetesAtracciones" pa
					INNER JOIN "atracciones_categorias" ac ON pa.id_atraccion = ac.id_atraccion
					INNER JOIN "subcategorias" sc ON ac.id_subcategoria = sc.id_subcategoria
					WHERE pa.id_paquete = pt.id_paquete_turistico
				),
				'[]'::jsonb
			),
			'hora_inicio', pt.hora_inicial,
			'atracciones', COALESCE(
				(
					SELECT jsonb_agg(
						jsonb_build_object(
							'id', at.id_atraccion,
							'nombre', at.nombre,
							'foto', COALESCE(
								(
									SELECT fa.foto
									FROM "fotos_atracciones" fa
									WHERE fa.id_atraccion = at.id_atraccion
									ORDER BY fa.orden ASC
									LIMIT 1
								),
								'N/A'
							),
							'duracion', pa.duracion,
							'descripcion', at.descripcion,
							'orden', pa.orden
						)
						ORDER BY pa.orden
					)
					FROM "PaquetesAtracciones" pa
					INNER JOIN "GestAtraccionesTuristicas" at ON pa.id_atraccion = at.id_atraccion
					WHERE pa.id_paquete = pt.id_paquete_turistico
				),
				'[]'::jsonb
			)
		) AS paquete
	FROM "GestPaquetesTuristicos" pt
	WHERE pt.id_paquete_turistico = ?`
