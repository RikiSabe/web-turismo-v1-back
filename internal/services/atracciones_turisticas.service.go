package services

var QueryAtraccionesTuristicasTODO = `
	SELECT 
		at.id_atracciones as id, at.tipo, at.nombre, at.ubicacion, at.descripcion, at.horarios, at.precio, at.estado
	FROM "GestAtraccionesTuristicas" as at
	ORDER BY id_atracciones ASC`

var QueryAtraccionesTuristicaUnique = `
	SELECT
		at.id_atracciones as id, at.tipo, at.nombre, at.ubicacion, at.descripcion, at.horarios, at.precio, at.estado
	FROM "GestAtraccionesTuristicas" as at
	WHERE id_atracciones = ?
	LIMIT 1`
