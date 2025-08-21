package services

var QueryAgenciasTODO = `
	SELECT 
		a.id_agencia as id, 
		a.nombre, 
		a.direccion, 
		a.telefono, 
		a.correo, 
		a.estado
	FROM "GestAgencias" as a
	ORDER BY id_agencia ASC`

var QueryAgenciaUnique = `
	SELECT
		ag.id_agencia as id, 
		ag.nombre, 
		ag.direccion, 
		ag.telefono, 
		ag.correo, 
		ag.estado,
		ag.descripcion,
		ag.id_departamento,
		ag.id_encargado,
	CASE
	WHEN COUNT(fa.id_foto) = 0 THEN to_jsonb('N/A'::text)
	ELSE jsonb_agg(
		DISTINCT jsonb_build_object(
			'id', fa.id_foto,
			'foto', fa.foto,
			'orden', fa.orden
		)
	)
	END AS fotos
	FROM "GestAgencias" as ag
	LEFT JOIN fotos_agencias fa ON fa.id_agencia = ag.id_agencia 
	WHERE ag.id_agencia = ?
	GROUP BY 
		ag.id_agencia, ag.nombre, ag.direccion, ag.telefono, 
		ag.correo, ag.estado, ag.descripcion, 
		ag.id_departamento, ag.id_encargado
	LIMIT 1;`
