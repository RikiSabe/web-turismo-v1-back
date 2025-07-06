package services

var QueryProvinciasTODO = `
	SELECT
		p.id_provincia as id,
		p.nombre
	FROM	
		provincias as p
	WHERE
		p.estado = True and p.id_departamento = ?`
