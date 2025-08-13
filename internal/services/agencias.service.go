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
			ag.id_departamento,
			ag.id_encargado
		FROM "GestAgencias" as ag
		WHERE ag.id_agencia = ?
		LIMIT 1`
