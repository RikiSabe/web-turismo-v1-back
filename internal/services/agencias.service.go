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
		a.id_agencia as id, 
		a.nombre, 
		a.direccion, 
		a.telefono, 
		a.correo, 
		a.estado
	FROM "GestAgencias" as a
	WHERE id_agencia = ?
	LIMIT 1`
